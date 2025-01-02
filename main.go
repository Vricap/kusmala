package main

import (
	"os"

	"github.com/vricap/kusmala/mode/file"
	"github.com/vricap/kusmala/mode/repl"
)

const DEV_MODE bool = true

func main() {
	arg := os.Args
	if len(arg) > 1 {
		file.Read(arg, DEV_MODE)
	} else {
		repl.Start(os.Stdin, os.Stdout, DEV_MODE)
	}
}
