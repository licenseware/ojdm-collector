#!/usr/bin/env bash
# Runs the acceptance suite on a real Windows host in Azure, driven entirely
# through the control plane: no RDP, no inbound NSG rules.
#
#   ./run.sh rc1          # label for this run's artifacts
#   ./run.sh --destroy    # drop the resource group when finished
#
# The collector and the suite are compiled here and uploaded, so the VM needs
# no Go toolchain. Results come back to ./out/results/.

set -euo pipefail

RG="${OJDM_RG:-ojdm-acceptance-rg}"
LOCATION="${OJDM_LOCATION:-westeurope}"
VM="${OJDM_VM:-ojdm-acceptance}"
VM_SIZE="${OJDM_VM_SIZE:-Standard_B2s}"
IMAGE="${OJDM_IMAGE:-MicrosoftWindowsServer:WindowsServer:2025-datacenter-azure-edition:latest}"
ADMIN_USER="${OJDM_ADMIN_USER:-ojdmadmin}"
CONTAINER=artifacts

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
OUT_DIR="$SCRIPT_DIR/out"
GO="${GO:-go}"

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

if ! command -v "$GO" >/dev/null 2>&1; then
  echo "go toolchain not found; install it or set GO=/path/to/go" >&2
  exit 1
fi

# Storage account names are global, lowercase, <=24 chars.
STORAGE="${OJDM_STORAGE:-ojdm$(az account show --query id -o tsv | tr -d '-' | cut -c1-16)}"

echo "==> building the collector and the suite for windows/amd64"
mkdir -p "$OUT_DIR"
( cd "$REPO_ROOT" && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 "$GO" build \
    -o "$OUT_DIR/ojdm-collector-$VARIANT.exe" . )
( cd "$REPO_ROOT" && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 "$GO" test -tags acceptance -c \
    -o "$OUT_DIR/acceptance-$VARIANT.test.exe" ./test/acceptance/ )

echo "==> ensuring resource group and storage"
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

upload() { # upload <local-file> <blob-name>
  az storage blob upload --account-name "$STORAGE" --account-key "$KEY" \
    --container-name "$CONTAINER" --name "$2" --file "$1" --overwrite --output none
}

echo "==> uploading artifacts"
upload "$OUT_DIR/ojdm-collector-$VARIANT.exe" "ojdm-collector-$VARIANT.exe"
upload "$OUT_DIR/acceptance-$VARIANT.test.exe" "acceptance-$VARIANT.test.exe"
upload "$SCRIPT_DIR/../plant_windows.ps1" "plant_windows.ps1"

if ! az vm show --resource-group "$RG" --name "$VM" --output none 2>/dev/null; then
  echo "==> creating VM $VM (no inbound rules; outbound only)"
  ADMIN_PASS="$(openssl rand -base64 24 | tr -d '/+=' | head -c 20)Aa1!"
  az vm create --resource-group "$RG" --name "$VM" --image "$IMAGE" --size "$VM_SIZE" \
    --admin-username "$ADMIN_USER" --admin-password "$ADMIN_PASS" \
    --nsg-rule NONE --public-ip-sku Standard --output none
  echo "    admin password (not needed for this run): $ADMIN_PASS"
fi

echo "==> running the acceptance suite on the VM (several minutes on a cold host)"
RENDERED="$OUT_DIR/run-suite-$VARIANT.ps1"
# A SAS token is a query string full of '&', which sed expands to the matched
# text in the replacement half. Escape it, or the VM gets a token-less URL and
# the storage account rejects the request as anonymous access.
SAS_ESCAPED="$(printf '%s' "$SAS_BASE" | sed -e 's/[&|\\]/\\&/g')"
sed -e "s|__ARTIFACT_SAS__|${SAS_ESCAPED}|" -e "s|__VARIANT__|${VARIANT}|" \
  "$SCRIPT_DIR/run-suite.ps1" > "$RENDERED"

RESULT="$(az vm run-command invoke --resource-group "$RG" --name "$VM" \
  --command-id RunPowerShellScript --scripts "@$RENDERED" \
  --query 'value[].message' -o tsv)"
echo "$RESULT" | sed 's/^/    /'

echo "==> downloading results"
az storage blob download-batch --account-name "$STORAGE" --account-key "$KEY" \
  --source "$CONTAINER" --pattern "results/$VARIANT-*" --destination "$OUT_DIR" --output none || true

echo
echo "run '$0 --destroy' when finished to drop the resource group."

# The suite's own exit code decides this script's, so it can gate anything.
if ! grep -q 'SUITEEXIT 0' <<<"$RESULT"; then
  echo "acceptance FAILED on windows" >&2
  exit 1
fi
echo "acceptance passed on windows"
