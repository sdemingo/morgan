package main

import (
	"strings"
	"unicode/utf8"
)

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
	}

	l.push(r)
	return textState
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
	return ch == ' ' || ch == '\t' || ch == '\n'
}
