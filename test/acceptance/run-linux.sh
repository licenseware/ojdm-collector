#!/usr/bin/env bash
# Runs the acceptance suite inside a Vagrant guest.
#
#   ./run-linux.sh            # every guest defined in the Vagrantfile
#   ./run-linux.sh debian     # one guest
#
# Compiles the collector and the suite on the host and uploads both, so the
# guest needs no Go toolchain and no synced folder.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
GO="${GO:-go}"
GUESTS=("$@")

if ! command -v "$GO" >/dev/null 2>&1; then
  echo "go toolchain not found; install it or set GO=/path/to/go" >&2
  exit 1
fi

if [[ ${#GUESTS[@]} -eq 0 ]]; then
  GUESTS=(debian rocky)
fi

STAGING="$(mktemp -d)"
trap 'rm -rf "$STAGING"' EXIT

# CGO_ENABLED=0 matches the release build. A cgo binary links against the
# building host's loader, which a guest will refuse to run.
echo "==> building the collector and the suite"
( cd "$REPO_ROOT" && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 "$GO" build -o "$STAGING/ojdm-collector" . )
( cd "$REPO_ROOT" && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 "$GO" test -tags acceptance -c \
    -o "$STAGING/acceptance.test" ./test/acceptance/ )

FAILED=()
for guest in "${GUESTS[@]}"; do
  echo "==> $guest"
  ( cd "$SCRIPT_DIR" && vagrant up "$guest" )

  ( cd "$SCRIPT_DIR"
    vagrant upload "$STAGING/ojdm-collector" /tmp/ojdm-collector "$guest"
    vagrant upload "$STAGING/acceptance.test" /tmp/acceptance.test "$guest"
    vagrant upload "$SCRIPT_DIR/plant_linux.sh" /tmp/plant_linux.sh "$guest" )

  if ( cd "$SCRIPT_DIR" && vagrant ssh "$guest" -c '
        set -e
        sudo chmod +x /tmp/ojdm-collector /tmp/acceptance.test /tmp/plant_linux.sh
        eval "$(sudo /tmp/plant_linux.sh)"
        sudo OJDM_BINARY=/tmp/ojdm-collector OJDM_JDK="$OJDM_JDK" \
          /tmp/acceptance.test -test.v' ); then
    echo "==> $guest PASSED"
  else
    echo "==> $guest FAILED"
    FAILED+=("$guest")
  fi
done

if [[ ${#FAILED[@]} -gt 0 ]]; then
  echo "acceptance failed on: ${FAILED[*]}" >&2
  exit 1
fi
echo "acceptance passed on: ${GUESTS[*]}"
