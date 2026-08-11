#!/usr/bin/env bash
# Provisioning only: makes sure a JDK exists on this host and prints the
# environment the acceptance suite needs. Everything else, including the
# planted layouts and every assertion, lives in the Go suite.
#
#   eval "$(sudo ./plant_linux.sh)"
#   sudo -E go test -tags acceptance ./test/acceptance/ -v
#
# Writes to system locations, so run it on a throwaway host only.

set -euo pipefail

if ! ls /usr/lib/jvm/*/bin/javac >/dev/null 2>&1; then
  if command -v apt-get >/dev/null; then
    export DEBIAN_FRONTEND=noninteractive
    apt-get update -qq >/dev/null
    apt-get install -y -qq default-jdk >/dev/null
  elif command -v dnf >/dev/null; then
    dnf install -y -q java-17-openjdk-devel >/dev/null
  else
    echo "no supported package manager; install a JDK manually" >&2
    exit 1
  fi
fi

JDK="$(readlink -f "$(dirname "$(dirname "$(ls /usr/lib/jvm/*/bin/javac | head -1)")")")"

echo "export OJDM_JDK=$JDK"
