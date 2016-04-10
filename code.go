package main

import (
	"fmt"
	"sync"
)

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

type Coder struct {
	lex    *Lexer
	wg     sync.WaitGroup
	output string
}

func HTMLCoder(l *Lexer) *Coder {
	g := &Coder{
		lex:    l,
		output: ""}

	g.wg.Add(1)
	go g.run()
	g.wg.Wait()

	return g
}

func (g *Coder) run() {
	g.output = ""
	for {
		tk, ok := <-g.lex.tokens
		if !ok {
			break
		}
		fmt.Printf("%d[%d][%s] ", tk.ttype, tk.offset, tk.value)
		switch tk.ttype {
		case header1Tk, header2Tk, header3Tk, header4Tk:
			codeHeader(g, tk, false)
			stack.Push(&tk)

		case newLineTk:
			tkContext := stack.Pop()
			if tkContext != nil && isHeader(tkContext) {
				codeHeader(g, *tkContext, true)
			}
			if tkContext != nil && tkContext.ttype == hyphenTk {
				codeList(g, *tkContext, true)
			}

		case textTk:
			g.output += tk.value

		case hyphenTk:
			codeList(g, tk, false)
			stack.Push(&tk)
		}

	}
	g.wg.Done()
}

func codeList(g *Coder, tk Token, close bool) {
	if close {
		g.output += "</li>"
	} else {
		g.output += "\n<li>"
	}
}

func codeHeader(g *Coder, tk Token, close bool) {
	switch tk.ttype {
	case header1Tk:
		if close {
			g.output += "</h1>\n"
		} else {
			g.output += "\n<h1>"
		}
	case header2Tk:
		if close {
			g.output += "</h2>\n"
		} else {
			g.output += "\n<h2>"
		}
	case header3Tk:
		if close {
			g.output += "</h3>\n"
		} else {
			g.output += "\n<h3>"
		}
	case header4Tk:
		if close {
			g.output += "</h4>\n"
		} else {
			g.output += "\n<h4>"
		}
	}
}
