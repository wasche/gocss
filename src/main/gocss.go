// gocss streaming CSS compressor
package main

import (
	"./lexer"
	"./parser"
	"fmt"
	"os"
	"bufio"
	"flag"
)

func process(lexer *lexer.Lexer, reader *bufio.Reader) {
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

func main() {
	parser := new(parser.Parser)
	lexer := new(lexer.Lexer)
	lexer.SetParser(parser)
	defer lexer.End()

	flag.Parse()

	if flag.NArg() == 0 {
		ri := os.Stdin
		reader := bufio.NewReader(ri)
		process(lexer, reader)
	}

	for i := 0; i < flag.NArg(); i++ {
		f, err := os.Open(flag.Arg(i))
		if f == nil {
			fmt.Fprintf(os.Stderr, "can't open %s: error %s\n", flag.Arg(1), err)
			os.Exit(1)
		}
		process(lexer, bufio.NewReader(f))
		f.Close()
	}
}
