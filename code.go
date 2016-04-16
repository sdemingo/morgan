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
	defer g.wg.Done()
	for {
		tk := g.next()
		if tk.ttype == nullTk {
			return
		}
		tkDispatcher(g, tk)
	}
}

func tkDispatcher(g *Coder, tk *Token) {
	fmt.Println(tk)

	checkFinishedLists(g, tk)

	switch tk.ttype {
	case header1Tk, header2Tk, header3Tk, header4Tk:
		codeHeader(g, tk)
	case italicTk, monoTk, ulineTk, boldTk:
		codeInline(g, tk)
	case textTk:
		g.output += strings.TrimSpace(tk.value) + " "
	case urlTk:
		codeDirectUrl(g, tk)
	case hyphenTk:
		codeItemList(g, tk)
	case newLineTk:
		codeNewLine(g, tk)
	}
}

func codeNewLine(g *Coder, tk *Token) {
	ntk := g.next()
	if ntk.ttype == newLineTk {
		g.output += "\n<br>\n"
	} else {
		g.back(ntk)
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
		ntk := g.next()
		if ntk.ttype == nullTk {
			return
		}
		if ntk.offset < itemOffset && ntk.ttype != newLineTk {
			g.back(ntk)
			break
		}
		tkDispatcher(g, ntk)
	}

	g.output += " </li>"
}

func codeInline(g *Coder, tk *Token) {
	g.output += " <" + tokenTag(tk) + "> "
	for {
		ntk := g.next()
		if ntk.ttype == nullTk {
			return
		}
		if ntk.ttype == tk.ttype {
			break
		}
		tkDispatcher(g, ntk)
	}
	g.output += " </" + tokenTag(tk) + "> "
}

func codeDirectUrl(g *Coder, tk *Token) {
	url := strings.TrimSpace(tk.value)

	ntk := g.next()
	if ntk.ttype == urlTextTk {
		g.output += "<a href=\"" + url + "\">" + ntk.value + "</a>"
	} else {
		g.output += "<a href=\"" + url + "\">" + url + "</a>"
		g.back(ntk)
	}
}

func codeHeader(g *Coder, tk *Token) {

	g.output += "\n<" + tokenTag(tk) + ">"
	for {
		ntk := g.next()
		if ntk.ttype == nullTk {
			return
		}
		if ntk.ttype == newLineTk {
			break
		}

		tkDispatcher(g, ntk)
	}
	g.output += "</" + tokenTag(tk) + ">\n"
}

func tokenTag(tk *Token) string {
	switch tk.ttype {
	case italicTk:
		return "i"
	case monoTk:
		return "code"
	case ulineTk:
		return "u"
	case boldTk:
		return "b"
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
