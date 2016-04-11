package main

import (
	"fmt"
	"strings"
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

func (g *Coder) getToken() *Token {
	tk, ok := <-g.lex.tokens
	if !ok {
		return nil
	}
	return &tk
}

func (g *Coder) run() {
	g.output = ""
	for {
		tk := g.getToken()
		if tk == nil {
			break
		}
		tkDispatcher(g, tk)

	}
	g.wg.Done()
}

func tkDispatcher(g *Coder, tk *Token) {
	if tk == nil {
		return
	}
	fmt.Printf("%d[%d][%s] ", tk.ttype, tk.offset, tk.value)
	switch tk.ttype {
	case header1Tk, header2Tk, header3Tk, header4Tk:
		codeHeader(g, tk)

	case italicTk:
		codeItalic(g, tk)
	case textTk:
		g.output += strings.TrimSpace(tk.value) + " "

	case hyphenTk:
		codeItemList(g, tk, false)
	}
}

func codeItemList(g *Coder, tk *Token, close bool) {
	g.output += "\n<li> "
	for {
		tk := g.getToken()
		if tk == nil {
			return
		}
		if tk.ttype == newLineTk {
			break
		}
		tkDispatcher(g, tk)
	}
	g.output += " </li>\n"
}

func codeItalic(g *Coder, tk *Token) {
	g.output += " <i> "
	for {
		tk := g.getToken()
		if tk == nil {
			return
		}
		if tk.ttype == italicTk {
			break
		}
		tkDispatcher(g, tk)
	}
	g.output += " </i> "
}

func codeHeader(g *Coder, tk *Token) {

	switch tk.ttype {
	case header1Tk:
		g.output += "\n<h1>"
	case header2Tk:
		g.output += "\n<h2>"
	case header3Tk:
		g.output += "\n<h3>"
	case header4Tk:
		g.output += "\n<h4>"
	}

	for {
		tk := g.getToken()
		if tk == nil {
			return
		}

		if tk.ttype == newLineTk {
			break
		}

		tkDispatcher(g, tk)
	}

	switch tk.ttype {
	case header1Tk:
		g.output += "</h1>\n"
	case header2Tk:
		g.output += "</h2>\n"
	case header3Tk:
		g.output += "</h3>\n"
	case header4Tk:
		g.output += "</h4>\n"
	}

}
