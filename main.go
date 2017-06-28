package main

import "fmt"

var version = "v0.1.0"
var dirty = ""

func main() {
	displayVersion := fmt.Sprintf("gack %s%s",
		version,
		dirty)
	Execute(displayVersion)
}
