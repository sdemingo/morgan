package morgan

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestHTMLCoder(t *testing.T) {
	content, err := ioutil.ReadFile("testfiles/prueba.org")
	if err != nil {
		panic(err)
	}
	s := string(content)

	parser := HTMLParser(s)
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println(parser.Output())
}
