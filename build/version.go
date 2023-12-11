package build

import (
	_ "embed"
	"strings"
)

var (
	Version string = strings.TrimSuffix(strings.TrimSuffix(strings.TrimSpace(version), "\n"), "\n")
	//go:embed version.txt
	version string
)
