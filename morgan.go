package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	content, err := ioutil.ReadFile("prueba.org")
	if err != nil {
		panic(err)
	}
	s := string(content)

	lexer := Lex(s)
	code := HTMLCoder(lexer)
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println(code.output)
}
