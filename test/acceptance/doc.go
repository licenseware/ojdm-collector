// Package acceptance verifies a built collector against real Java
// installations on a real host: it plants the installation layouts and failure
// modes the collector has to cope with, scans them, and asserts the report,
// the debug log and its rotation.
//
// Every test is guarded by the "acceptance" build tag, because the suite
// installs software and writes to system locations. Run it with:
//
//	go test -tags acceptance ./test/acceptance/
//
// See README.md for the environment it expects and for the drivers that run it
// inside a throwaway guest.
//
// This file carries no build tag on purpose: without it the package holds no
// files for a default build, and `go vet ./...` and `go test ./...` fail with
// "build constraints exclude all Go files".
package acceptance
