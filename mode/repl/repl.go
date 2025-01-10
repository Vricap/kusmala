package repl

import (
	"bufio"
	"fmt"
	"io"
	"os/user"

	"github.com/vricap/kusmala/lexer"
	"github.com/vricap/kusmala/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer, DEV_MODE bool) {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Halo %s! Ini adalah bahasa pemrograman KUSMALA!\n", user.Username)
	fmt.Println("Silahkan untuk mengetik perintah.")
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

		if len(pars.DevErrors) != 0 && DEV_MODE {
			printDevError(pars.DevErrors, out)
			continue
		}

		if len(pars.Errors) != 0 {
			printParsingError(pars.Errors, out)
			continue
		}

		parser.PrintTree(tree.Statements)
	}
}

func printParsingError(err []string, out io.Writer) {
	fmt.Println("Pesan error mungkin tidak akurat :)")
	for _, e := range err {
		io.WriteString(out, "\t"+e+"\n")
	}
}

func printDevError(err []string, out io.Writer) {
	for _, e := range err {
		io.WriteString(out, "\t"+e+"\n")
	}
}
