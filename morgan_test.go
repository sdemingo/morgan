package morgan

import (
	"fmt"
	"testing"
)

func TestHTMLHeaders(t *testing.T) {
	test := `
* Test Article title
** TODO Test Subheader to do
** DONE Test Subheader done
*** Test Subheader
   DEADLINE: <2016-05-10 mar 13:15>

Meto un timestamp <2016-05-10 mar> estándar para ver como se exporta

* Other header
`
	parser := HTMLParser(test)
	html := parser.Output()

	fmt.Println(html)
}

func TestHTMLInlines(t *testing.T) {
	test := `Lorem ipsum /dolor sit amet/, consectetuer adipiscing
*elit*. Donec hendrerit tempor tellus. =Donec pretium posuere=
tellus. Proin quam nisl, tincidunt et, mattis eget, convallis nec,
purus.

Lorem ipsum /dolor sit amet/.
`
	parser := HTMLParser(test)
	html := parser.Output()
	fmt.Println(html)
}

func TestHTMLLinks(t *testing.T) {

	test := `
Lorem ipsum [[http://www.url.com][ipsum dolor sit amet]],
consectetuer adipiscing elit. Donec hendrerit tempor tellus. Donec
pretium posuere tellus. [[http://www.url.com]] Proin quam nisl,
tincidunt et, mattis eget, convallis nec, purus.

[[/dir/images/files/image.jpg]]
   `

	parser := HTMLParser(test)
	html := parser.Output()
	fmt.Println(html)
}

func TestHTMLItemLists(t *testing.T) {
	test := `
Lorem ipsum dolor sit amet, - bluff lista consectetuer adipiscing:
  - Item 1
  - Item 2
  - Proin quam nisl, tincidunt et.
  - Item 3

Otro tipo de lista sin margen:

- Item 1
- Item 2


`

	// Añado lista enlazada:
	//    - Item 1
	//    - Nueva lista:
	//      - Subitem 1
	//      - Subitem 2
	//    - Item 3

	// `
	// Ahora pongo una segunda lista ordenada:
	//    1. Item 1
	//    2. Item 2
	//    3. Item 3
	// CuCu
	// `
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

// Pellentesque dapibus. Preparamos  una lista.

// `
// 	parser := HTMLParser(test)
// 	html := parser.Output()
// 	fmt.Println(html)
// }
