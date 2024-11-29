package repl

import (
	"bufio"
	"fmt"
	"io"
)

const PROMPT = ">>"

func Start(input io.Reader, output io.Writer) {
	bufio.NewScanner(input)

	for {
		fmt.Printf(PROMPT)
	}
}
