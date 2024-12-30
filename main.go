package main

import (
	"os"

	"github.com/vricap/kusmala/mode/file"
	"github.com/vricap/kusmala/mode/repl"
)

func main() {
	arg := os.Args
	if len(arg) > 1 {
		file.Read(arg)
	} else {
		repl.Start(os.Stdin, os.Stdout)
	}
}
