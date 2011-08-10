// Minification parser
package parser

import "./lexer"
import "fmt"
import "os"

type Parser struct {
}

func (p *Parser) Token(token lexer.Token, value string) {
	os.Stdout.WriteString(value)
	//fmt.Fprintf(os.Stderr, "token: %s, value: %s\n", token, value)
}

func (p *Parser) End() {
	fmt.Fprintf(os.Stderr, "Done\n")
}
