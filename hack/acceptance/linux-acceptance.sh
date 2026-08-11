#!/usr/bin/env bash
# Acceptance suite for a release candidate, run as root inside a throwaway
# Linux guest. Installs a real JDK, plants the layouts and failure modes each
# fix was written for, then asserts the observable contract: the report, the
# debug log beside it, and the rotation of a previous log.
#
#   usage: linux-acceptance.sh /path/to/ojdm-collector

set -uo pipefail

BINARY="${1:?usage: linux-acceptance.sh /path/to/ojdm-collector}"
WORK=/opt/ojdm-acceptance
PASS=0
FAIL=0

check() { # check <description> <condition-result>
  if [[ "$2" == "0" ]]; then
    echo "PASS $1"
    PASS=$((PASS + 1))
  else
    echo "FAIL $1"
    FAIL=$((FAIL + 1))
  fi
}

# ------------------------------------------------------------------ java setup
if ! ls /usr/lib/jvm/*/bin/javac >/dev/null 2>&1; then
  if command -v apt-get >/dev/null; then
    export DEBIAN_FRONTEND=noninteractive
    apt-get update -qq >/dev/null && apt-get install -y -qq default-jdk >/dev/null
  else
    dnf install -y -q java-17-openjdk-devel >/dev/null
  fi
fi
JDK="$(readlink -f "$(dirname "$(dirname "$(ls /usr/lib/jvm/*/bin/javac | head -1)")")")"
echo "INFO jdk=$JDK"

# ------------------------------------------------------- planted failure modes
rm -rf "$WORK"; mkdir -p "$WORK"

# 1. An installation under a root whose own name contains "bin".
mkdir -p /opt/bin-tools
rm -rf /opt/bin-tools/jdk-planted
cp -rL "$JDK" /opt/bin-tools/jdk-planted || echo "PLANT-ERROR cp failed"
check "planted install is in place" "$([[ -x /opt/bin-tools/jdk-planted/bin/java ]] && echo 0 || echo 1)"

# 2. A directory that cannot be listed: must warn, must not stop the scan.
rm -rf /opt/aaa_denied; mkdir -p /opt/aaa_denied/bin; chmod 000 /opt/aaa_denied

# 3. A path too long to stat: a non-permission error, the one that used to
#    abort the whole search root. Built relatively so creation itself succeeds.
rm -rf /opt/aaa_broken; mkdir -p /opt/aaa_broken
( cd /opt/aaa_broken || exit 0
  segment="$(printf 'd%.0s' {1..200})"
  for _ in $(seq 1 30); do mkdir -p "$segment" && cd "$segment" || break; done )

# ------------------------------------------------------------------ run 1
REPORT="$WORK/report.csv"
LOG="$WORK/logs/debug.log"
"$BINARY" -search-paths=/opt/bin-tools,/opt/aaa_denied,/opt/aaa_broken \
  -output-path="$REPORT" >"$WORK/console1.txt" 2>&1
check "run 1 exits successfully" "$?"

check "csv report created" "$([[ -s $REPORT ]] && echo 0 || echo 1)"
check "debug log created beside the report" "$([[ -s $LOG ]] && echo 0 || echo 1)"

python3 - "$REPORT" "$LOG" <<'PY'
import csv, json, os, sys
report, logpath = sys.argv[1], sys.argv[2]

rows = list(csv.DictReader(open(report)))
entries = [json.loads(l) for l in open(logpath) if l.strip()]
results = []

results.append(("report contains the planted bin-named install",
                any(r["JavaHome"].startswith("/opt/bin-tools/jdk-planted") for r in rows)))
results.append(("every row resolves its vm shared library",
                bool(rows) and all(r["DynLibBinPath"] and os.path.exists(r["DynLibBinPath"]) for r in rows)))
results.append(("every jdk row resolves javac",
                all(os.path.exists(r["JavaCBinPath"]) for r in rows if r["IsJDK"] == "true")))
results.append(("hostname recorded on every row", all(r["HostName"] for r in rows)))

results.append(("debug log is valid json throughout", len(entries) > 0))
results.append(("log records the search paths",
                any(e["message"] == "scanning for java installations" for e in entries)))
results.append(("non-permission failure reported and survived",
                any(e.get("permission_denied") is False and "aaa_broken" in e.get("path", "") for e in entries)))
results.append(("debug detail present in the file",
                any(e["level"] == "debug" for e in entries)))
results.append(("final entry points the operator at the log",
                entries[-1]["message"].startswith("done")))

for description, ok in results:
    print(("PASS " if ok else "FAIL ") + description)
sys.exit(1 if any(not ok for _, ok in results) else 0)
PY
check "report and log assertions all hold" "$?"
PYLINE="$(python3 - "$REPORT" "$LOG" <<'PY'
import csv, json, sys
print(len(list(csv.DictReader(open(sys.argv[1])))), len([l for l in open(sys.argv[2]) if l.strip()]))
PY
)"
echo "INFO rows_and_log_entries=$PYLINE"

check "console stays quiet at info level" \
  "$(grep -q 'found java file' "$WORK/console1.txt" && echo 1 || echo 0)"

# ------------------------------------------------------------------ run 2
sleep 1
"$BINARY" -search-paths=/opt/bin-tools -output-path="$REPORT" >"$WORK/console2.txt" 2>&1
ARCHIVES=$(find "$WORK/logs" -name 'debug-*.log' | wc -l)
check "previous log archived on rerun" "$([[ $ARCHIVES -eq 1 ]] && echo 0 || echo 1)"
check "archive kept the earlier run's content" \
  "$(grep -lq 'aaa_denied' "$WORK"/logs/debug-*.log && echo 0 || echo 1)"
check "fresh log started for the new run" \
  "$(grep -q 'aaa_denied' "$LOG" && echo 1 || echo 0)"

# ------------------------------------------------------------------ run 3
"$BINARY" -search-paths=/opt/bin-tools -output-path="$WORK/r3.csv" \
  -log-path="$WORK/custom/my.log" -log-level=debug >"$WORK/console3.txt" 2>&1
check "-log-path honoured" "$([[ -s $WORK/custom/my.log ]] && echo 0 || echo 1)"
check "-log-level=debug surfaces detail on the console" \
  "$(grep -q 'found java file' "$WORK/console3.txt" && echo 0 || echo 1)"

# ------------------------------------------------------------------ run 4
# The permission branch is only reachable unprivileged: root reads everything.
UNPRIV="$(ls /home | head -1)"
chown -R "$UNPRIV" "$WORK" 2>/dev/null
su "$UNPRIV" -c "$BINARY -search-paths=/opt/aaa_denied,/opt/bin-tools -output-path=$WORK/r4.csv \
  -log-path=$WORK/unpriv/debug.log" >"$WORK/console4.txt" 2>&1
check "unprivileged run still produces a report" "$([[ -s $WORK/r4.csv ]] && echo 0 || echo 1)"
check "permission failure reported with its flag" \
  "$(grep -q '"permission_denied":true' "$WORK/unpriv/debug.log" && echo 0 || echo 1)"
check "scan continued past the permission failure" \
  "$(grep -q 'jdk-planted' "$WORK/r4.csv" && echo 0 || echo 1)"

chmod 755 /opt/aaa_denied 2>/dev/null
echo "SUITE pass=$PASS fail=$FAIL"

# The exit code is what gates CI; without it a failed assertion reads as success.
[[ $FAIL -eq 0 ]] || exit 1
