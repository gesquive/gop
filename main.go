package main

import "fmt"

var version = "v0.2.2"
var dirty = ""

func main() {
	displayVersion := fmt.Sprintf("gop %s%s",
		version,
		dirty)
	Execute(displayVersion)
}
