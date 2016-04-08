package main

import (
	"strings"
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
	newLineTk

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

func initState(l *Lexer) stateFunc {

	r := l.next()
	if r == eof || !utf8.ValidRune(r) {
		return nil
	}

	if isWhitespace(r) {
		return consume
	}

	switch r {
	case '*':
		return headerState
	case '\n':
		return newLineState
	}

	l.push(r)
	return textState
}

func newLineState(l *Lexer) stateFunc {
	var r rune
	for {
		r = l.next()
		if r == eof {
			return nil
		}

		if r != '\n' {
			break
		}
	}

	l.emit(newLineTk)
	l.push(r)
	return initState
}

func headerState(l *Lexer) stateFunc {

	if strings.HasPrefix(l.input[l.pos:], "*") {
		l.pos++
		l.emit(header2Tk)
		return initState
	}
	if strings.HasPrefix(l.input[l.pos:], "**") {
		l.pos = l.pos + 2
		l.emit(header3Tk)
		return initState
	}
	if strings.HasPrefix(l.input[l.pos:], "***") {
		l.pos = l.pos + 3
		l.emit(header4Tk)
		return initState
	}

	l.emit(header1Tk)
	return initState
}

func textState(l *Lexer) stateFunc {
	var r rune
	for {
		r = l.next()
		if r == eof {
			return nil
		}

		if isWhitespace(r) {
			break
		}
	}
	l.emit(textTk)
	l.push(r)
	return initState
}

func consume(l *Lexer) stateFunc {
	var r rune
	for {
		r = l.next()
		if r == eof {
			return nil
		}

		if !isWhitespace(r) {
			break
		}
	}
	l.push(r)
	return initState
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}
