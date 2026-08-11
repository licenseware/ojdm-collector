# Acceptance suite

Verifies a built collector against **real Java installations on a real host**:
it plants the installation layouts and failure modes the collector has to cope
with, scans them, and asserts the report, the debug log and its rotation.

The unit tests answer "is the logic right". This answers "does it work on a
machine that looks like a customer's".

> [!WARNING]
> The suite installs packages and writes to system locations
> (`/opt`, `/Applications`, `C:\Program Files\Java`, `C:\bin-tools`). Run it on
> a throwaway guest, never on a workstation. That is why it sits behind the `acceptance`
> build tag and never runs during `go test ./...`.

## Layout

| Path | Purpose |
|---|---|
| `acceptance_test.go` | Every assertion, for every platform |
| `support_linux_test.go`, `support_darwin_test.go`, `support_windows_test.go` | Platform specifics: where to plant, how to make something unreadable |
| `plant_linux.sh`, `plant_darwin.sh`, `plant_windows.ps1` | Provisioning only: make a JDK exist, print where it is |
| `Vagrantfile` | Throwaway Debian and Rocky guests |
| `run-linux.sh` | Builds, uploads and runs the suite in those guests |
| `azure/run.sh`, `azure/run-suite.ps1` | Same, on a Windows Server VM in Azure |

The assertions live in Go so there is exactly one copy of them. The scripts
only provision hosts.

## Running it

### Linux, via Vagrant

```bash
cd test/acceptance
vagrant up
./run-linux.sh              # both guests
./run-linux.sh debian       # just one
```

Two distributions on purpose: the deb and rpm families lay their JDKs out
differently. The collector and the suite are compiled on your machine and
uploaded, so the guests need no Go toolchain.

### Windows, via Azure

```bash
az login && az account set --subscription <id>

cd test/acceptance/azure
./run.sh rc1                # label for this run's artifacts
./run.sh --destroy          # when finished
```

Runs entirely through the Azure control plane: no RDP, no inbound firewall
rules. It creates the resource group, storage account and VM on first use and
reuses the VM afterwards. Results land in `azure/out/results/`.

**Delete the resource group when you are done** — the VM bills until you do.

### Directly on a host you do not mind breaking

```bash
# Linux, as root
eval "$(sudo ./plant_linux.sh)"
sudo -E env "PATH=$PATH" OJDM_BINARY=/path/to/ojdm-collector OJDM_JDK="$OJDM_JDK" \
  go test -tags acceptance -v ./test/acceptance/
```

```bash
# macOS, as a normal admin user — not root, or the permission test skips
eval "$(./plant_darwin.sh)"
OJDM_BINARY=/path/to/ojdm-collector go test -tags acceptance -v ./test/acceptance/
```

```powershell
# Windows, as Administrator
.\plant_windows.ps1 | Invoke-Expression
$env:OJDM_BINARY = 'C:\path\to\ojdm-collector.exe'
go test -tags acceptance -v .\test\acceptance\
```

### In CI

`.github/workflows/ci.yml` runs the whole thing on every pull request and on
every push to `main`, on GitHub's Linux, Windows and macOS runners. The branch
ruleset requires those checks, so nothing reaches `main` without them; the
release build itself does not re-run them.

## Environment

| Variable | Required | Meaning |
|---|---|---|
| `OJDM_BINARY` | yes | Collector under test |
| `OJDM_JDK` | yes | A JDK holding `javac`, `jps` and `jinfo`; emitted by the planting script |
| `OJDM_JRE` | no | A JRE-only installation, so the JDK and JRE paths are both covered |
| `OJDM_EXTRA_ROOT` | no | A second drive or mount point; installations there are outside the default search paths |

Optional variables are skipped with a log line rather than failing, so the
suite still runs on a host without a second drive.

## What it asserts

- Every planted installation appears, including one under a root whose own name
  contains `bin`, and one on a second drive
- No installation is reported twice, however many search paths reach it
- Every row resolves its VM shared library, and every JDK row resolves `javac`
- A running JVM is identified through `jps` and `jinfo`
- An entry that cannot be read is reported **and the scan continues** — the
  defect that used to produce blank reports on customer hosts
- A permission failure is flagged as one (Linux drops privileges first, since
  the job runs as root and root reads everything; macOS runs unprivileged
  throughout)
- The debug log sits beside the report, is JSON throughout, keeps `debug`
  detail the console does not show, and ends by pointing at itself
- `-log-path` and `-log-level` behave
- A rerun archives the previous log instead of destroying it

## Extending it

Add an assertion as an ordinary Go test in `acceptance_test.go`. Use
`newScanFixture(t)` if you need the planted host and a completed scan; it
returns the parsed rows and log entries. It plants once per test, so a test
that only needs a plain run should call `runCollector` directly.

Anything platform specific — a location to plant in, a way to make the
filesystem misbehave — belongs behind a helper in the `support_*_test.go`
files, so the assertions themselves stay platform neutral.

If a new assertion needs software on the host, install it in the planting
script and export its location, rather than installing from Go: that keeps the
suite runnable against a host somebody else provisioned.
