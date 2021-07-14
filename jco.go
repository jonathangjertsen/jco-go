package main

import (
	"github.com/jonathangjertsen/jco-go/buildconfig"
	"github.com/jonathangjertsen/jco-go/cmd"
)

func main() {
	defer buildconfig.PanicHandler()
	cmd.Execute()
}
