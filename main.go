package main

import (
	"fmt"

	"github.com/vricap/kusmala/lexer"
	"github.com/vricap/kusmala/parser"
)

func main() {
	// user, err := user.Current()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Halo %s! Ini adalah bahasa pemrograman KUSMALA!\n", user.Username)
	// fmt.Println("Silahkan untuk mengetik perintah.")
	// repl.Start(os.Stdin)

	input := `buat x = 1 + 1;`
	lex := lexer.NewLex(input)
	pars := parser.NewPars(lex)
	ast := pars.ParsCode()

	for i := 0; i < len(ast.Statements); i++ {
		fmt.Println(ast.Statements[i])
	}
}
