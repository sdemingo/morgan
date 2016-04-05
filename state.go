package main

import "strings"

const (
	headerMark  = "*"
	newLineMark = "\n"
)

func lexHeader(l *Lexer) stateFunc {
	for {
		if strings.HasPrefix(l.input[l.pos:], headerMark) {
			break
		}
		if l.next() == eof {
			break
		}
	}

	l.start = l.pos
	l.pos = strings.Index(l.input[l.pos:], newLineMark)
	if l.pos < 0 {
		return nil
	}
	l.emit(headerTk)

	return nil
}
