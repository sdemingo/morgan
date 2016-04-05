package main

import (
	"fmt"
	"io/ioutil"
	"unicode/utf8"
)

type stateFunc func(*Lexer) stateFunc

var initialState = lexHeader

type Token struct {
	ttype int
	value string
}

const (
	eof = iota
	headerTk
)

type Lexer struct {
	input  string
	start  int // start position of this item
	pos    int // current position in the input
	width  int
	tokens chan Token
}

func Lex(input string) *Lexer {
	l := &Lexer{
		input:  input,
		tokens: make(chan Token),
	}
	go l.run() // Concurrently run state machine.
	return l
}

func (l *Lexer) run() {
	for state := initialState; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *Lexer) emit(ttype int) {
	l.tokens <- Token{ttype, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *Lexer) next() (rune rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	rune, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return rune
}

func main() {
	content, err := ioutil.ReadFile("prueba.org")
	if err != nil {
		panic(err)
	}
	s := string(content)

	lex := Lex(s)
	for {
		tk := <-lex.tokens
		fmt.Printf("%v", tk)
	}
}
