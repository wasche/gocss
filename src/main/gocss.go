// gocss streaming CSS compressor
package main

import (
	"./lexer"
	"./parser"
	"fmt"
	"os"
	"flag"
	"runtime"
	"strings"
)

// general options
var suffixGenerated *string = flag.String("g", "-gen.css", "Suffix of generated files")
var suffixCompressed *string = flag.String("c", "-c.css", "Suffix of compressed files")
var verbose *bool = flag.Bool("v", false, "Print progress information")
var yui *bool = flag.Bool("y", false, "Match output to YUI Compressor v2.4.6")
// right to left conversion
var convert *bool = flag.Bool("r", false, "Convert for right to left languages")
var convertGen *bool = flag.Bool("R", false, "Convert generated file for right to left languages")
var suffixRTLS *string = flag.String("G", "-rtl.css", "Suffix of generated RTL files")
var suffixRTL *string = flag.String("C", "-rtl-c.css", "Suffix of compressed RTL files")
// configuration file
var config *string = flag.String("i", "Makefile.gcs", "File to read configuration from")
var createConfig *bool = flag.Bool("T", false, "Output a sample configuration file")
var generate *bool = flag.Bool("o", false, "Output generated files")

func process(in *os.File, out *os.File) {
	// set up channels
	runes := make(chan(int))
	tokenValues := make(chan(lexer.TokenValue))
	minified := make(chan(string))
	eof := make(chan(int))

	// TODO multi file input streamer
	ifs := &InputFileStreamer{In: in, Out: runes}
	lexer := &lexer.Lexer{In: runes, Out: tokenValues}
	// TODO generated file using TokenValueFileStreamer
	// TODO generated RTL - dual channel, pass-through and output
	parser := &parser.Parser{In: tokenValues, Out: minified, Yui: *yui}
	// TODO compressed RTL - dual channel, pass-through and output
	ofs := &OutputFileStreamer{In: minified, Out: out, Eof: eof}

	go ifs.Run()
	go lexer.Run()
	go parser.Run()
	go ofs.Run()

	// wait for chain to finish
	<- eof
}

func processFile(name string, i int) {
	target := strings.Replace(name, *suffixGenerated, *suffixCompressed, 1)
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
