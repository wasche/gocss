// gocss streaming CSS compressor
package main

import (
	"./lexer"
	"./parser"
	"./rtl"
	"fmt"
	"os"
	"flag"
	"runtime"
	"strings"
)

// general options
var stdin *bool = flag.Bool("i", false, "Read from <STDIN> and write compressed data to <STDOUT>")
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
var config *string = flag.String("f", "Makefile.gcs", "File to read configuration from")
var createConfig *bool = flag.Bool("T", false, "Output a sample configuration file")
var generate *bool = flag.Bool("o", false, "Output generated files")

func main() {
	flag.Parse()

	// std in, of option selected
	if *stdin {
		stream()
		return
	}

	// list of files given on command line
	if flag.NArg() > 0 {
		convertArgs()
		return
	}

	// read from configuration file
	cfg, err := os.Open(*config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open config file (%s): %s\n", *config, err)
		os.Exit(2)
	}
	defer cfg.Close()
}

// convert from stdin
func stream() {
	// set up channels
	runes := make(chan(int))
	tokenValues := make(chan(lexer.TokenValue))
	minified := make(chan(string))
	eof := make(chan(int))

	ifs := &InputFileStreamer{In: os.Stdin, Out: runes}
	lexer := &lexer.Lexer{In: runes, Out: tokenValues}
	parser := &parser.Parser{In: tokenValues, Out: minified, Yui: *yui}
	ofs := &OutputFileStreamer{In: minified, Out: os.Stdout, Eof: eof}

	go ifs.Run()
	go lexer.Run()
	go parser.Run()
	go ofs.Run()

	// wait for chain to finish
	<- eof
}

// convert list of files given on command line
func convertArgs() {
	threads := 4
	n := flag.NArg()
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

func processFile(name string, threadNum int) {
	fi, err := os.Open(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open %s: error %s\n", name, err)
		return
	}
	defer fi.Close()

	target := strings.Replace(name, *suffixGenerated, *suffixCompressed, 1)
	fo, err := os.Create(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't create %s: error %s\n", target, err)
		return
	}
	defer fo.Close()

	ifs, runes := CreateInputFileStreamer(fi)
	go ifs.Run()

	lexer, tokenValues := lexer.CreateLexer(runes)
	go lexer.Run()

	if *convertGen {
		rtlGenName := strings.Replace(name, *suffixGenerated, *suffixRTLS, 1)
		rtlGenFile, err := os.Create(rtlGenName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't create %s: error %s\n", rtlGenName, err)
			return
		}
		defer rtlGenFile.Close()

		if *verbose { fmt.Fprintf(os.Stderr, "[%d] Converting: %s\n", threadNum, rtlGenName) }

		// split channels
		gensplitter, out1, out2 := CreateChannelSplitter(tokenValues)
		go gensplitter.Run()
		genconverter := rtl.CreateConverter(out2, rtlGenFile)
		go genconverter.Run()
		tokenValues = out1
	}

	if *verbose { fmt.Fprintf(os.Stderr, "[%d] Compressing: %s\n", threadNum, target) }

	parser, minified := parser.CreateParser(tokenValues, *yui)
	go parser.Run()

	if *convert {
		rtlName := strings.Replace(name, *suffixGenerated, *suffixRTL, 1)
		rtlFile, err := os.Create(rtlName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't create %s: error %s\n", rtl, err)
			return
		}
		defer rtlFile.Close()

		if *verbose { fmt.Fprintf(os.Stderr, "[%d] Converting: %s\n", threadNum, rtlName) }

		// split channels
		splitter, out3, out4 := CreateChannelSplitter(minified)
		go splitter.Run()
		converter := rtl.CreateConverter(out4, rtlFile)
		go converter.Run()
		minified = out3
	}

	ofs, eof := CreateOutputFileStreamer(minified, fo)
	go ofs.Run()

	// wait for chain to finish
	<- eof
}
