package main

import (
	gl "github.com/kubex-ecosystem/logz"
)

func main() {

	c := gl.GetLoggerZ("")
	c.Info("Hello, World!")

	MMin()
}
