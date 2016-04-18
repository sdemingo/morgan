package morgan

import (
	"fmt"
	"testing"
)

func TestHTMLHeaders(t *testing.T) {
	test := `* Test Article title
`
	parser := HTMLParser(test)
	html := parser.Output()

	fmt.Println(html)
}

// func TestHTMLCoder(t *testing.T) {
// 	content, err := ioutil.ReadFile("testfiles/prueba.org")
// 	if err != nil {
// 		panic(err)
// 	}
// 	s := string(content)

// 	parser := HTMLParser(s)
// 	fmt.Println()
// 	fmt.Println()
// 	fmt.Println()
// 	fmt.Println(parser.Output())
// }
