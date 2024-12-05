package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/vricap/kusmala/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Halo %s! Ini adalah bahasa pemrograman KUSMALA!\n", user.Username)
	fmt.Println("Silahkan untuk mengetik perintah.")
	repl.Start(os.Stdin)

	// input := `buat x = 1 + 1 - 1;
	// buat x = 1 + 1;`
	// lex := lexer.NewLex(input)
	// pars := parser.NewPars(lex)
	// ast := pars.ParsCode()

	// for i := 0; i < len(ast.Statements); i++ {
	// 	fmt.Println(ast.Statements[i].TokenLiteral())
	// }
}
