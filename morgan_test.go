package morgan

import (
	"fmt"
	"testing"
)

func TestHTMLHeaders(t *testing.T) {
	test := `* Test Article title
** Test Subheader
*** Test Subheader
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
`
	parser := HTMLParser(test)
	html := parser.Output()
	fmt.Println(html)
}

func TestHTMLLinks(t *testing.T) {

	test := `Lorem [[http://www.url.com][ipsum dolor sit amet]],
consectetuer adipiscing elit. Donec hendrerit tempor tellus. Donec
pretium posuere tellus. [[http://www.url.com]] Proin quam nisl,
tincidunt et, mattis eget, convallis nec, purus.   `

	parser := HTMLParser(test)
	html := parser.Output()
	fmt.Println(html)
}
