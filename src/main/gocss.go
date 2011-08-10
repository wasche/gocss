// gocss streaming CSS compressor
package main

import (
	"./lexer"
	"./parser"
	"fmt"
	"os"
	"bufio"
)

func main() {
	parser := new(parser.Parser)
	lexer := new(lexer.Lexer)
	lexer.SetParser(parser)
	defer lexer.End()

	ri := os.Stdin
	reader := bufio.NewReader(ri)

	for {
		switch r, s, er := reader.ReadRune(); {
		case s < 0:
			fmt.Fprintf(os.Stderr, "error reading: %s\n", er.String())
			os.Exit(1)
		case s == 0: // EOF
			return
		case s > 0:
			lexer.Tokenize(r)
		}
	}
}
