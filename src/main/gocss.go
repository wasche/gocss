// gocss streaming CSS compressor
package main

import (
	"./lexer"
	"./parser"
	"fmt"
	"os"
	"bufio"
	"flag"
	"runtime"
)

func process(f *os.File) {
	parser := new(parser.Parser)
	lexer := new(lexer.Lexer)
	lexer.SetParser(parser)
	defer lexer.End()
	reader := bufio.NewReader(f)
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

func processFiles(name string, result chan(int)) {
	f, err := os.Open(name)
	if f == nil {
		fmt.Fprintf(os.Stderr, "can't open %s: error %s\n", name, err)
		os.Exit(1)
	}
	process(f)
	f.Close()
	result <- 0
}

func main() {
	flag.Parse()
	n := flag.NArg()

	if flag.NArg() == 0 {
		process(os.Stdin)
		return
	}
	
	threads := 4
	runtime.GOMAXPROCS(threads)

	result := make(chan int, threads)
	for i := 0; i < n; i++ {
		go processFiles(flag.Arg(i), result)
	}

	// wait for all jobs to complete
	for i := 0; i < n; i++ {
		<-result
	}
}
