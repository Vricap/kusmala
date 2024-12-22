package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/vricap/kusmala/lexer"
	"github.com/vricap/kusmala/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		input := scanner.Text()
		lex := lexer.NewLex(input)
		pars := parser.NewPars(lex)
		tree := pars.ConstructTree()

		if len(pars.Errors) != 0 {
			printError(pars.Errors, out)
			continue
		}

		parser.PrintTree(tree.Statements)
	}
}

func printError(err []string, out io.Writer) {
	for _, e := range err {
		io.WriteString(out, "\t"+e+"\n")
	}
}
