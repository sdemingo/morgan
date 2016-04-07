package main

import "unicode/utf8"

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
	l.emit(headerTk)
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
