package main

import (
	"AlekseyPromet/examples/simplewebhook/commands"
)

var (
	xBuildVersion    string
	xBuildHashCommit string
)

func main() {
	println("xBuildVersion", xBuildVersion)
	println("xBuildHashCommit", xBuildHashCommit)

	commands.Execute()
}
