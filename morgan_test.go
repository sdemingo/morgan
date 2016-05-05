package morgan

import (
	"fmt"
	"testing"
)

// func TestHTMLHeaders(t *testing.T) {
// 	test := `
// * Test Article title
// ** Test Subheader
// *** Test Subheader
// * Other header
// `
// 	parser := HTMLParser(test)
// 	html := parser.Output()

// 	fmt.Println(html)
// }

// func TestHTMLInlines(t *testing.T) {
// 	test := `Lorem ipsum /dolor sit amet/, consectetuer adipiscing
// *elit*. Donec hendrerit tempor tellus. =Donec pretium posuere=
// tellus. Proin quam nisl, tincidunt et, mattis eget, convallis nec,
// purus.

// Lorem ipsum /dolor sit amet/.
// `
// 	parser := HTMLParser(test)
// 	html := parser.Output()
// 	fmt.Println(html)
// }

// func TestHTMLLinks(t *testing.T) {

// 	test := `
// Lorem ipsum [[http://www.url.com][ipsum dolor sit amet]],
// consectetuer adipiscing elit. Donec hendrerit tempor tellus. Donec
// pretium posuere tellus. [[http://www.url.com]] Proin quam nisl,
// tincidunt et, mattis eget, convallis nec, purus.

// [[/dir/images/files/image.jpg]]
//    `

// 	parser := HTMLParser(test)
// 	html := parser.Output()
// 	fmt.Println(html)
// }

func TestHTMLItemLists(t *testing.T) {
	test := `
Lorem ipsum dolor sit amet, - bluff lista consectetuer adipiscing:
  - Item 1
  - Item 2:
  - Proin quam nisl, tincidunt et.
  - Nueva lista

    bla bla bla
  - Item 3

Bla bla
`
	parser := HTMLParser(test)
	html := parser.Output()
	fmt.Println(html)
}

// func TestHTMLCode(t *testing.T) {
// 	test := `

// #+TITLE: Test title
// #+PROPERTY: prop1

// :PROPERTIES:
// :Title:    Goldberg Variations
// :Composer: J.S. Bach
// :END:

// Pellentesque dapibus. Preparamos  una lista.

// Lorem ipsum dolor sit amet, consectetuer adipiscing:
// #+BEGIN_SRC sh
// go get github.com/sdemingo/morgan
// rm -rf /foo/bar
// #+END_SRC

// Bla bla

// `
// 	parser := HTMLParser(test)
// 	html := parser.Output()
// 	fmt.Println(html)
// }
