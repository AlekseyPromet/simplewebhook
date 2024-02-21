package main

import (
	"AlekseyPromet/examples/simplewebhook/commands"
)

var (
	xBuildVersion    string
	xBuildHashCommit string
)

func main() {
	if xBuildVersion != "" {
		println("xBuildVersion", xBuildVersion)
		println("xBuildHashCommit", xBuildHashCommit)
	}

	commands.Execute()
}
