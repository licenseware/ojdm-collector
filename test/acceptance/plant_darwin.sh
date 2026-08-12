#!/usr/bin/env bash
# Provisioning only: makes sure a JDK exists on this host and prints the
# environment the acceptance suite needs. Everything else, including the
# planted layouts and every assertion, lives in the Go suite.
#
#   eval "$(./plant_darwin.sh)"
#   go test -tags acceptance ./test/acceptance/ -v
#
# The suite plants into /Applications, so run it on a throwaway Mac only.
# No sudo: /Applications is writable by the admin group, and the permission
# test only means something when the run is unprivileged.

set -euo pipefail

# java_home is authoritative on macOS and already satisfied by the JDKs
# preinstalled on the CI images.
find_jdk() {
  local home
  home="$(/usr/libexec/java_home 2>/dev/null || true)"
  if [ -n "$home" ] && [ -x "$home/bin/javac" ] && [ -x "$home/bin/jps" ]; then
    echo "$home"
  fi
}

JDK="$(find_jdk)"

if [ -z "$JDK" ]; then
  if command -v brew >/dev/null; then
    brew install --cask temurin >/dev/null
    JDK="$(find_jdk)"
  fi
fi

if [ -z "$JDK" ]; then
  echo "no jdk with javac and jps found; install one manually" >&2
  exit 1
fi

echo "export OJDM_JDK=$JDK"
