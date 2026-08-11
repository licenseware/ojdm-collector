#!/usr/bin/env bash
# End-to-end Windows reproduction for the "blank CSV" collector bug, driven
# entirely through the Azure control plane: no RDP, no inbound NSG rules.
#
#   ./run-azure-repro.sh buggy    # before the fix  -> expect VERDICT=REPRODUCED
#   ./run-azure-repro.sh fixed    # after the fix   -> expect VERDICT=NOT_REPRODUCED
#   ./run-azure-repro.sh --destroy
#
# Whatever is in the working tree right now is what gets built and shipped; the
# argument is only the label used for artifact and result filenames.

set -euo pipefail

RG="${OJDM_RG:-ojdm-repro-rg}"
LOCATION="${OJDM_LOCATION:-westeurope}"
VM="${OJDM_VM:-ojdm-repro}"
VM_SIZE="${OJDM_VM_SIZE:-Standard_B2s}"
IMAGE="${OJDM_IMAGE:-MicrosoftWindowsServer:WindowsServer:2025-datacenter-azure-edition:latest}"
ADMIN_USER="${OJDM_ADMIN_USER:-ojdmadmin}"
CONTAINER=artifacts

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
OUT_DIR="$SCRIPT_DIR/out"
GO="${GO:-$HOME/.nix-profile/bin/go}"

if [[ "${1:-}" == "--destroy" ]]; then
  echo "==> deleting resource group $RG"
  az group delete --name "$RG" --yes --no-wait
  exit 0
fi

VARIANT="${1:-}"
if [[ -z "$VARIANT" ]]; then
  echo "usage: $0 <variant-label|--destroy>" >&2
  exit 2
fi

# Storage account names are global, lowercase, <=24 chars.
STORAGE="${OJDM_STORAGE:-ojdm$(az account show --query id -o tsv | tr -d '-' | cut -c1-16)}"

echo "==> building windows/amd64 collector (label: $VARIANT)"
mkdir -p "$OUT_DIR"
(cd "$REPO_ROOT" && GOOS=windows GOARCH=amd64 "$GO" build -o "$OUT_DIR/ojdm-$VARIANT.exe" .)

echo "==> ensuring resource group + storage"
az group create --name "$RG" --location "$LOCATION" --output none
az storage account create --name "$STORAGE" --resource-group "$RG" --location "$LOCATION" \
  --sku Standard_LRS --kind StorageV2 --min-tls-version TLS1_2 --output none
KEY="$(az storage account keys list --account-name "$STORAGE" --resource-group "$RG" \
  --query '[0].value' -o tsv)"
az storage container create --name "$CONTAINER" --account-name "$STORAGE" \
  --account-key "$KEY" --output none

# Expiry is deliberately short: the SAS is embedded in the script the VM runs.
EXPIRY="$(date -u -d '+4 hours' '+%Y-%m-%dT%H:%MZ')"
SAS="$(az storage container generate-sas --name "$CONTAINER" --account-name "$STORAGE" \
  --account-key "$KEY" --permissions rwl --expiry "$EXPIRY" --https-only -o tsv)"
SAS_BASE="https://${STORAGE}.blob.core.windows.net/${CONTAINER}?${SAS}"

echo "==> uploading collector binary and the acceptance suite"
az storage blob upload --account-name "$STORAGE" --account-key "$KEY" \
  --container-name "$CONTAINER" --name "ojdm-$VARIANT.exe" \
  --file "$OUT_DIR/ojdm-$VARIANT.exe" --overwrite --output none
# The VM-side scripts fetch the suite by name; keep it in step with the repo.
az storage blob upload --account-name "$STORAGE" --account-key "$KEY" \
  --container-name "$CONTAINER" --name "windows-acceptance.ps1" \
  --file "$REPO_ROOT/hack/acceptance/windows-acceptance.ps1" --overwrite --output none

if ! az vm show --resource-group "$RG" --name "$VM" --output none 2>/dev/null; then
  echo "==> creating VM $VM (no inbound rules; outbound only)"
  ADMIN_PASS="$(openssl rand -base64 24 | tr -d '/+=' | head -c 20)Aa1!"
  az vm create --resource-group "$RG" --name "$VM" --image "$IMAGE" --size "$VM_SIZE" \
    --admin-username "$ADMIN_USER" --admin-password "$ADMIN_PASS" \
    --nsg-rule NONE --public-ip-sku Standard --output none
  echo "    admin password (not needed for this run): $ADMIN_PASS"
fi

echo "==> running repro on the VM (several minutes; downloads a JRE + 3 scans)"
PS_SCRIPT="${OJDM_SCRIPT:-$SCRIPT_DIR/repro.ps1}"
RENDERED="$OUT_DIR/$(basename "${PS_SCRIPT%.ps1}")-$VARIANT.ps1"
# A SAS token is a query string full of '&', which sed expands to the matched
# text in the replacement half. Escape it, or the VM gets a token-less URL and
# the storage account rejects the request as anonymous access.
SAS_ESCAPED="$(printf '%s' "$SAS_BASE" | sed -e 's/[&|\\]/\\&/g')"
sed -e "s|__ARTIFACT_SAS__|${SAS_ESCAPED}|" -e "s|__VARIANT__|${VARIANT}|" \
  "$PS_SCRIPT" > "$RENDERED"

az vm run-command invoke --resource-group "$RG" --name "$VM" \
  --command-id RunPowerShellScript --scripts "@$RENDERED" \
  --query 'value[].message' -o tsv | sed 's/^/    /'

echo "==> downloading result CSVs and logs"
az storage blob download-batch --account-name "$STORAGE" --account-key "$KEY" \
  --source "$CONTAINER" --pattern "results/$VARIANT-*" --destination "$OUT_DIR" --output none
find "$OUT_DIR/results" -name "$VARIANT-*" -printf '    %p (%s bytes)\n' 2>/dev/null || true

echo
echo "run '$0 --destroy' when done to drop the resource group."
