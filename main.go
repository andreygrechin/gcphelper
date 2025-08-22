package main

import (
	"os"

	"github.com/andreygrechin/gcphelper/cmd"
)

var (
	Version   string // Version is set during build time.
	BuildTime string // BuildTime is set during build time.
	Commit    string // Commit is set during build time.
)

func main() {
	exitCode := cmd.Execute(cmd.VersionInfo{
		Version:   Version,
		BuildTime: BuildTime,
		Commit:    Commit,
	})
	os.Exit(exitCode)
}
