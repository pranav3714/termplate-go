package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "dev"
	Commit    = "unknown"
	Date      = "unknown"
	Branch    = "unknown"
	GoVersion = runtime.Version()
)

type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	Branch    string `json:"branch"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
}

func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		Branch:    Branch,
		GoVersion: GoVersion,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func (i Info) String() string {
	return fmt.Sprintf("%s (commit: %s, built: %s, %s)",
		i.Version, i.Commit, i.Date, i.GoVersion)
}
