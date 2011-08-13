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

var verbose *bool = flag.Bool("v", false, "Print progress information")

func process(f *os.File) {
	parser := new(parser.Parser)
	lexer := new(lexer.Lexer)
	lexer.SetParser(parser)
	defer lexer.End()
	reader := bufio.NewReader(f)
	for {
		r, s, er := reader.ReadRune()
		switch {
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

func processFiles(i int, queue chan(string), result chan(int)) {
	for {
		select {
		case name := <-queue:
			f, err := os.Open(name)
			if f == nil {
				fmt.Fprintf(os.Stderr, "can't open %s: error %s\n", name, err)
				os.Exit(1)
			}
			if *verbose { fmt.Fprintf(os.Stderr, "[%d] Compressing %s\n", i, name) }
			process(f)
			f.Close()
		default:
			result <- 0
			return
		}
	}
}

func main() {
	flag.Parse()
	n := flag.NArg()

	if flag.NArg() == 0 {
		process(os.Stdin)
		return
	}
	
	threads := 4

	if threads > n { threads = n }

	runtime.GOMAXPROCS(threads)

	queue := make(chan string, 64)
	result := make(chan int, threads)
	for i := 0; i < threads; i++ {
		go processFiles(i, queue, result)
	}

	for i := 0; i < n; i++ {
		queue <- flag.Arg(i)
	}

	// wait for all jobs to complete
	for i := 0; i < threads; i++ {
		<-result
	}
}
