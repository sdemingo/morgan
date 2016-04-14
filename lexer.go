package main

import (
	"fmt"
	"strings"
	"sync"
	"unicode/utf8"
)

type stateFunc func(*Lexer) stateFunc

type Token struct {
	ttype  int
	value  string
	offset int
}

func (tk *Token) String() string {
	return fmt.Sprintf("[%d (%d)(%s)] ", tk.ttype, tk.offset, tk.value)
}

var eof = rune(0)
var nullToken = &Token{ttype: nullTk, value: "NullToken"}

const (
	nullTk    = iota
	header1Tk //1
	header2Tk
	header3Tk
	header4Tk
	textTk // 5
	newLineTk
	hyphenTk
	ulistTk
	italicTk
	monoTk //10
	ulineTk
	boldTk
	urlTk
)

type Lexer struct {
	input  string
	start  int // start position of this item
	pos    int // current position in the input
	width  int
	offset int // token offset in the line
	tokens chan *Token
	wg     sync.WaitGroup
}

func Lex(input string) *Lexer {
	l := &Lexer{
		input:  input,
		tokens: make(chan *Token, 5),
		offset: 0,
	}

	go l.lex()
	return l
}

func (l *Lexer) lex() {
	for state := lexInitState; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *Lexer) emit(ttype int) {
	l.tokens <- &Token{ttype, l.input[l.start:l.pos], l.offset}
	l.start = l.pos
}

func (l *Lexer) read() (rune rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	rune, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return rune
}

func (l *Lexer) prev() (rune rune) {
	rune, _ = utf8.DecodeRuneInString(l.input[l.pos-2 : l.pos-1])
	return rune
}

func (l *Lexer) unread(rune rune) {
	l.pos -= utf8.RuneLen(rune)
	if l.pos < 0 {
		l.pos = 0
	}
}

func (l *Lexer) incOffset(rune rune) {
	if isWhitespace(rune) || rune == '-' {
		l.offset++
	}

	if isNewline(rune) {
		l.offset = 0
	}
}

func lexInitState(l *Lexer) stateFunc {

	r := l.read()
	if r == eof || !utf8.ValidRune(r) {
		return nil
	}

	l.incOffset(r)

	if isWhitespace(r) {
		return consume
	}

	if isPunctuation(r) {
		return lexInitState
	}

	switch r {
	case '*':
		prev := l.prev()
		next := l.read()
		if isNewline(prev) && isWhitespace(next) {
			l.unread(next)
			return headerState
		}
		l.unread(next)
		l.emit(boldTk)
		return lexInitState
	case '\n':
		l.offset = 0
		return newLineState
	case '-':
		l.emit(hyphenTk)
		return lexInitState
	case '/':
		l.emit(italicTk)
		return lexInitState
	case '_':
		l.emit(ulineTk)
		return lexInitState
	case '=':
		l.emit(monoTk)
		return lexInitState
	}

	l.unread(r)
	return textState
}

func newLineState(l *Lexer) stateFunc {
	var r rune
	for {
		r = l.read()
		if r == eof {
			return nil
		}

		if r != '\n' {
			break
		}
	}
	l.unread(r)
	l.emit(newLineTk)
	return lexInitState
}

func headerState(l *Lexer) stateFunc {

	if strings.HasPrefix(l.input[l.pos:], "* ") {
		l.pos++
		l.emit(header2Tk)
		return lexInitState
	}
	if strings.HasPrefix(l.input[l.pos:], "** ") {
		l.pos = l.pos + 2
		l.emit(header3Tk)
		return lexInitState
	}
	if strings.HasPrefix(l.input[l.pos:], "*** ") {
		l.pos = l.pos + 3
		l.emit(header4Tk)
		return lexInitState
	}

	l.emit(header1Tk)
	return lexInitState
}

func textState(l *Lexer) stateFunc {
	var r rune

	if strings.HasPrefix(l.input[l.pos:], "http://") {
		return urlState
	}

	for {
		// consume runes of a word until spaces, reserved
		// characters or  punctuaction
		r = l.read()
		if r == eof {
			return nil
		}

		if isEndWord(r) {
			break
		}

		if r == '/' || r == '-' || r == '=' || r == '_' || r == '*' || r == '=' {
			nr := l.read()
			if nr != eof && isEndWord(nr) {
				l.unread(nr)
				break
			}
		}
	}

	l.unread(r)
	l.emit(textTk)
	return lexInitState
}

func urlState(l *Lexer) stateFunc {
	var r rune
	for {
		// consume runes of a word until spaces
		r = l.read()
		if r == eof {
			return nil
		}

		if isWhitespace(r) || isNewline(r) {
			break
		}
	}
	l.unread(r)
	l.emit(urlTk)
	return lexInitState
}

func consume(l *Lexer) stateFunc {
	var r rune
	for {
		r = l.read()
		if r == eof {
			return nil
		}

		if !isWhitespace(r) || !isNewline(r) {
			break
		}
	}
	l.unread(r)
	return lexInitState
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

func isNewline(ch rune) bool {
	return ch == '\n'
}

func isEndWord(ch rune) bool {
	return isWhitespace(ch) || isNewline(ch) || isPunctuation(ch)
}

func isPunctuation(ch rune) bool {
	return ch == '.'
}

func isHeader(tk *Token) bool {
	return tk.ttype == header1Tk || tk.ttype == header2Tk || tk.ttype == header3Tk || tk.ttype == header4Tk
}
