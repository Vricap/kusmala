package file

import (
	"fmt"
	"log"
	"os"

	"github.com/vricap/kusmala/ast"
	"github.com/vricap/kusmala/evaluator"
	"github.com/vricap/kusmala/lexer"
	"github.com/vricap/kusmala/object"
	"github.com/vricap/kusmala/parser"
)

func Read(arg []string, DEV_MODE bool) {
	tree := readFile(arg[1], DEV_MODE)
	env := object.NewEnv()
	evals := evaluator.Eval(tree, env)
	if evals != nil {
		printEval(evals)
	}
	if len(arg) > 2 {
		switch arg[2] {
		case "-tree":
			parser.PrintTree(tree.Statements)
		default:
			log.Fatal("Argumen tidak diketahui!")
		}
	}
}

func readFile(path string, DEV_MODE bool) *ast.Tree {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Mengalami masalah saat membuka file tersebut:\n%s", err)
	}
	if !isKusmalaFile(path) {
		log.Fatal("File merupakan bukan file kusmala. File kusmala ektensi '.km'")
	}

	lex := lexer.NewLex(string(data))
	pars := parser.NewPars(lex)
	tree := pars.ConstructTree()
	if len(pars.DevErrors) != 0 && DEV_MODE {
		printDevError(pars.DevErrors)
	}
	if len(pars.Errors) != 0 {
		printParsingError(pars.Errors)
	}
	return tree
}

func isKusmalaFile(n string) bool {
	x := n[len(n)-3 : len(n)]
	return x == ".km"
}

func printParsingError(err []string) {
	fmt.Println("Pesan error mungkin tidak akurat :)")
	for _, e := range err {
		fmt.Println("\t" + e)
	}
	os.Exit(1)
}

func printDevError(err []string) {
	for _, e := range err {
		fmt.Println("\t" + e)
	}
	os.Exit(1)
}

func printEval(evals []object.Object) {
	// for _, eval := range evals {
	// 	fmt.Printf("%s\n", eval.Inspect())
	// }
}
