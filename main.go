package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome! %s to the Moneky Programming Language. \n", user.Username)
	fmt.Println("Try out some monkey commands to get started")

	repl.Start(os.Stdin, os.Stdout)
}
