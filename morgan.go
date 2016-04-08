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

	lex := Lex(s)
	parser := Build(lex.tokens)
	fmt.Printf("\n%s\n", parser.code)
}
