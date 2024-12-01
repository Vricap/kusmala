package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/vricap/kusmala/lexer"
	"github.com/vricap/kusmala/token"
)

const PROMPT = ">>"

func Start(input io.Reader) {
	scanner := bufio.NewScanner(input)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		lex := lexer.NewLex(line)

		tok := lex.NextToken()
		for tok.Type != token.EOF {
			fmt.Printf("%+v\n", tok)
			tok = lex.NextToken()
		}
	}
}
