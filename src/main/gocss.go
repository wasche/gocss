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
	"strings"
)

var verbose *bool = flag.Bool("v", false, "Print progress information")
var regexFrom *string = flag.String("f", ".css", "String to replace")
var regexTo *string = flag.String("t", "-c.css", "String to replace with (default: -c.css)")
var yui *bool = flag.Bool("y", false, "Match output to YUI Compressor v2.4.6")
// right to left conversion
var convert *bool = flag.Bool("c", false, "Convert for right to left languages")
var regexRTL *string = flag.String("r", "-rtl-c.css", "String to replace with for RTL (default: -rtl-c.css)")
var convertSource *bool = flag.Bool("C", false, "Convert source file for right to left languages")
var regexRTLS *string = flag.String("R", "-rtl.css", "String to replace with for RTL source (default: -rtl.css)")
// configuration file
var config *string = flag.String("i", "gocss.cfg", "File to read configuration from (default: gocss.cfg)")

func process(in *os.File, out *os.File) {
	parser := &parser.Parser{Output: out, Yui: *yui}
	lexer := &lexer.Lexer{Parser: parser}
	defer lexer.End()
	reader := bufio.NewReader(in)
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

func processFile(name string, i int) {
	target := strings.Replace(name, *regexFrom, *regexTo, 1)
	fi, err := os.Open(name)
	fo, err := os.Create(target)
	defer fi.Close()
	defer fo.Close()
	if fi == nil {
		fmt.Fprintf(os.Stderr, "can't open %s: error %s\n", name, err)
		os.Exit(1)
	}
	if *verbose { fmt.Fprintf(os.Stderr, "[%d] Compressing %s\n", i, name) }
	process(fi, fo)
}

func processFiles(i int, queue chan(string), result chan(int)) {
	for {
		select {
		case name := <-queue:
			processFile(name, i)
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
		// check for configuration file

		process(os.Stdin, os.Stdout)
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
	for i := 0 ; i < threads; i++ {
		<-result
	}
}
