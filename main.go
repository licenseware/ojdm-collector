// Command ojdm-collector finds the JVMs and JDKs installed on a host, reports
// their versions and the java processes running out of them.
package main

import "github.com/licenseware/ojdm-collector/internal/cmd"

// Stamped at build time by GoReleaser via -ldflags -X main.<name>=...
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	cmd.Run(cmd.Build{
		Version: version,
		Commit:  commit,
		Date:    date,
		BuiltBy: builtBy,
	})
}
