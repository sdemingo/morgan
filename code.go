package main

import "fmt"

type Stack []*Token

func (s *Stack) Push(v *Token) {
	*s = append(*s, v)
}

func (s *Stack) Pop() *Token {
	if len(*s) <= 0 {
		return nil
	}
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}

var stack Stack

func (l *Lexer) code() {
	l.output = ""
	for {
		tk, ok := <-l.tokens
		if !ok {
			break
		}
		fmt.Printf("%d[%s] ", tk.ttype, tk.value)
		switch tk.ttype {
		case header1Tk, header2Tk, header3Tk, header4Tk:
			codeHeader(l, tk, false)
			stack.Push(&tk)

		case newLineTk:
			tkContext := stack.Pop()
			if tkContext != nil && isHeader(tkContext) {
				codeHeader(l, *tkContext, true)
			}

		case textTk:
			l.output += tk.value
		}
	}
	l.wg.Done()
}

func codeHeader(l *Lexer, tk Token, close bool) {
	switch tk.ttype {
	case header1Tk:
		if close {
			l.output += "</h1>\n"
		} else {
			l.output += "\n<h1>"
		}
	case header2Tk:
		if close {
			l.output += "</h2>\n"
		} else {
			l.output += "\n<h2>"
		}
	case header3Tk:
		if close {
			l.output += "</h3>\n"
		} else {
			l.output += "\n<h3>"
		}
	case header4Tk:
		if close {
			l.output += "</h4>\n"
		} else {
			l.output += "\n<h4>"
		}
	}
}
