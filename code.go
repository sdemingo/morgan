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
		return nullToken
	}
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}

func (s *Stack) Top() *Token {
	if len(*s) <= 0 {
		return nullToken
	}
	return (*s)[len(*s)-1]
}

type Coder struct {
	lex    *Lexer
	wg     sync.WaitGroup
	stack  Stack
	backed Stack
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

func (g *Coder) next() *Token {
	tk := g.backed.Pop()
	if tk.ttype != nullTk {
		return tk
	}

	tk, ok := <-g.lex.tokens
	if ok {
		return tk
	}
	return nullToken
}

func (g *Coder) back(t *Token) {
	g.backed.Push(t)
}

func (g *Coder) run() {
	g.output = ""
	for {
		tk := g.next()
		if tk.ttype == nullTk {
			break
		}
		tkDispatcher(g, tk)

	}
	g.wg.Done()
}

func tkDispatcher(g *Coder, tk *Token) {

	fmt.Println(tk)

	checkFinishedLists(g, tk)

	switch tk.ttype {
	case header1Tk, header2Tk, header3Tk, header4Tk:
		codeHeader(g, tk)
	case italicTk:
		codeItalic(g, tk)
	case textTk:
		g.output += strings.TrimSpace(tk.value) + " "
	case hyphenTk:
		codeItemList(g, tk)
	}
}

func codeItemList(g *Coder, tk *Token) {

	itemOffset := tk.offset + 1
	rootListToken := Token{ulistTk, "ul", itemOffset - 1}

	if g.stack.Top().ttype != ulistTk {
		g.stack.Push(&rootListToken)
		g.output += "\n<ul>"
	}

	g.output += "\n<li> "

	for {
		tk := g.next()
		if tk.offset < itemOffset && tk.ttype != newLineTk {
			g.back(tk)
			break
		}
		tkDispatcher(g, tk)
	}

	g.output += " </li>"
}

func codeItalic(g *Coder, tk *Token) {
	g.output += " <i> "
	for {
		tk := g.next()
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
		tk := g.next()
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

func checkFinishedLists(g *Coder, tk *Token) {

	tkContext := g.stack.Top()
	if tkContext.ttype == ulistTk {
		if tk.ttype != newLineTk && tk.offset < tkContext.offset {
			g.output += "\n</ul>\n"
			g.stack.Pop()
		}

		// TODO: in org mode lists can be finished with a
		// double newLine character

	}
}
