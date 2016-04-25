package morgan

import (
	"fmt"
	"strings"
	"sync"
)

const debugMode = true

type Stack []*Token

func (s *Stack) push(v *Token) {
	*s = append(*s, v)
}

func (s *Stack) pop() *Token {
	if len(*s) <= 0 {
		return &Token{nullTk, "", 0}
	}
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}

func (s *Stack) top() *Token {
	if len(*s) <= 0 {
		return &Token{nullTk, "", 0}
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

func HTMLParser(input string) *Coder {
	l := newLexer(input)
	g := &Coder{
		lex:    l,
		output: ""}

	g.wg.Add(1)
	go g.run()
	g.wg.Wait()

	return g
}

func (g *Coder) Output() string {
	return g.output
}

func (g *Coder) next() *Token {
	tk := g.backed.pop()
	if tk.ttype != nullTk {
		return tk
	}

	tk, ok := <-g.lex.tokens
	if ok {
		return tk
	}
	return &Token{nullTk, "", 0}
}

func (g *Coder) back(t *Token) {
	g.backed.push(t)
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
	if debugMode {
		fmt.Println(tk)
	}

	closeOpenedLists(g, tk)

	if g.stack.top().ttype == codeTk {
		if tk.ttype != codeTk {
			g.output += tk.value
			return
		}
	}

	switch tk.ttype {
	case header1Tk, header2Tk, header3Tk, header4Tk:
		codeHeader(g, tk)
	case italicTk, monoTk, ulineTk, boldTk:
		codeInline(g, tk)
	case textTk:
		codeText(g, tk)
	case blankTk:
		g.output += " "
	case urlTk:
		codeUrl(g, tk)
	case hyphenTk:
		codeItemList(g, tk)
	case newLineTk:
		codeNewLine(g, tk)
	case codeTk:
		codeSnippets(g, tk)
	case propBlockTk:
		g.output += "<!--\n" + tk.value + "\n-->"
	}
}

func codeText(g *Coder, tk *Token) {
	stk := g.stack.top()
	if !isHeader(stk) && stk.ttype != ulistTk {
		if stk.ttype != parTk {
			g.stack.push(&Token{parTk, "", 0})
			g.output += "\n<p>"
		}
	}

	g.output += tk.value
}

func codeNewLine(g *Coder, tk *Token) {

	// close tags which must be closed with one breakline
	closeOpenedHeader(g, tk)

	ntk := g.next()
	// close tags which must be closed with two breakline or a
	// breakline and a new container opentag
	if ntk.ttype == newLineTk || ntk.ttype == nullTk || ntk.ttype == codeTk {
		closeOpenedPar(g, tk)
	}
	g.back(ntk)

	g.output += " "

}

func codeSnippets(g *Coder, tk *Token) {

	if g.stack.top().ttype == codeTk {
		g.output += "\n</code></pre>\n"
		g.stack.pop()
		return
	}

	g.output += "\n<pre><code>\n"
	g.stack.push(tk)
}

func codeItemList(g *Coder, tk *Token) {

	itemOffset := tk.offset
	rootListToken := Token{ulistTk, "ul", itemOffset}

	if g.stack.top().ttype != ulistTk {
		// before opened a list, close other container tags as
		// paragraphs, ...
		if g.stack.top().ttype == parTk {
			g.stack.pop()
			g.output += "</p>\n"
		}
		g.stack.push(&rootListToken)
		g.output += "\n<ul>"
	}

	g.output += "\n<li>"

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
	g.output += "<" + tokenTag(tk) + ">"
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
	g.output += "</" + tokenTag(tk) + ">"
}

func codeUrl(g *Coder, tk *Token) {
	url := strings.TrimSpace(tk.value)

	ntk := g.next()
	if ntk.ttype == urlTextTk {
		g.output += "<a href=\"" + url + "\">" + ntk.value + "</a>"
	} else {
		if isImageUrl(url) {
			g.output += "<img src=\"" + url + "\"/>"
		} else {
			g.output += "<a href=\"" + url + "\">" + url + "</a>"
		}
		g.back(ntk)
	}
}

func codeHeader(g *Coder, tk *Token) {
	g.output += "\n<" + tokenTag(tk) + ">"
	g.stack.push(tk)
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
	case parTk:
		return "p"
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

func closeOpenedLists(g *Coder, tk *Token) {

	tkContext := g.stack.top()
	if tkContext.ttype == ulistTk {

		if tk.ttype != newLineTk && tk.ttype != blankTk && tk.offset < tkContext.offset {
			g.output += "\n</ul>\n"
			g.stack.pop()
		}

		// TODO: in org mode lists can be finished with a
		// double newLine character
	}
}

func closeOpenedPar(g *Coder, tk *Token) {
	if tk.ttype != newLineTk {
		return
	}

	if g.stack.top().ttype == parTk {
		g.output += "</p>\n"
		g.stack.pop()
	} //  else {
	// 	g.output += "\n<br>\n"
	// }
}

func closeOpenedHeader(g *Coder, tk *Token) {
	if tk.ttype != newLineTk {
		return
	}
	stk := g.stack.top()
	if isHeader(stk) { //tags which must be closed with one breakline
		g.output += "</" + tokenTag(stk) + ">\n"
		g.stack.pop()
		return
	}
}

func isImageUrl(url string) bool {
	return strings.HasSuffix(url, ".jpg") || strings.HasSuffix(url, ".png")
}
