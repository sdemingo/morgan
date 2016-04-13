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
	case italicTk, monoTk:
		codeInline(g, tk)
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

func codeInline(g *Coder, tk *Token) {
	g.output += " <" + tokenTag(tk) + "> "
	for {
		ntk := g.next()
		if ntk.ttype == tk.ttype {
			break
		}
		tkDispatcher(g, ntk)
	}
	g.output += " </" + tokenTag(tk) + "> "
}

func codeHeader(g *Coder, tk *Token) {

	g.output += "\n<" + tokenTag(tk) + ">"
	for {
		tk := g.next()
		if tk.ttype == newLineTk {
			break
		}

		tkDispatcher(g, tk)
	}
	g.output += "</" + tokenTag(tk) + ">\n"
}

func tokenTag(tk *Token) string {
	switch tk.ttype {
	case italicTk:
		return "i"
	case monoTk:
		return "code"
	case header1Tk:
		return "h1"
	case header2Tk:
		return "h2"
	case header3Tk:
		return "h3"
	case header4Tk:
		return "h4"
	}
	return ""
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
