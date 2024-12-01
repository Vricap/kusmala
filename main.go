package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/vricap/kusmala/repl"
)

func main() {
	// in := `1 + 1`
	// lex := lexer.NewLex(in)

	// tok := lex.NextToken()
	// for tok.Type != token.EOF {
	// 	// fmt.Printf("%v\n", tok)
	// 	tok = lex.NextToken()
	// }

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Halo %s! Ini adalah bahasa pemrograman KUSMALA!\n", user.Username)
	fmt.Println("Silahkan untuk mengetik perintah.")
	repl.Start(os.Stdin)
}
