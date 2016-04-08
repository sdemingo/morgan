package main

import "fmt"

type Generator struct {
	tokens chan Token // input tokens from lexer
	output chan string
	code   string
}

func Build(input chan Token) *Generator {
	gn := &Generator{input, make(chan string), ""}
	go gn.run()
	gn.code = <-gn.output
	return gn
}

func (g *Generator) run() {
	for {
		tk, ok := <-g.tokens
		if !ok {
			break
		}
		fmt.Printf("%d ", tk.ttype)

	}
	g.output <- "fin"
}
