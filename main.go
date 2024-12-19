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
	repl.Start(os.Stdin, os.Stdout)
}
