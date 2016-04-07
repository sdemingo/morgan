package main

import (
	"fmt"
	"io/ioutil"
	"unicode/utf8"
)

type stateFunc func(*Lexer) stateFunc

type Token struct {
	ttype int
	value string
}

var eof = rune(0)

const (
	header1Tk = iota + 1
	header2Tk
	header3Tk
	header4Tk
	textTk

	headerPreffix   = "*"
	preCodePreffix  = "="
	italicPreffix   = "/"
	listItemPreffix = "-"
	linkPreffix     = "[["
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
	for state := initState; state != nil; {
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

func (l *Lexer) push(rune rune) {
	l.pos -= utf8.RuneLen(rune)
	if l.pos < 0 {
		l.pos = 0
	}
}

func main() {
	content, err := ioutil.ReadFile("prueba.org")
	if err != nil {
		panic(err)
	}
	s := string(content)

	lex := Lex(s)

	for {
		tk, ok := <-lex.tokens
		fmt.Printf("%d ", tk.ttype)
		if !ok {
			break
		}
	}
}
