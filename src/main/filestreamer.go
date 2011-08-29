/**
 * @author wasche
 * @since 2011.08.28
 */
package main

import (
	"os"
	"bufio"
	"fmt"
	"./lexer"
)

type InputFileStreamer struct {
	In  *os.File
	Out chan(int)
}

func (fs *InputFileStreamer) Run() {
	reader := bufio.NewReader(fs.In)
	for {
		r, s, er := reader.ReadRune()
		switch {
		case s < 0:
			fmt.Fprintf(os.Stderr, "error reading: %s\n", er.String())
			os.Exit(1)
		case s == 0: // EOF
			fs.Out <- -1
			return
		case s > 0:
			fs.Out <- r
		}
	}
}

var ZERO_STR string

type OutputFileStreamer struct {
	In  chan(string)
	Out *os.File
	Eof chan(int)
}

func (fs *OutputFileStreamer) Run() {
	var s string
	for {
		s = <- fs.In
		switch (s) {
		case ZERO_STR:
			if fs.Eof != nil { fs.Eof <- 0 }
			return
		default:
			fs.Out.WriteString(s)
		}
	}
}

type TeeFileStreamer struct {
	In   chan(string)
	Out  chan(string)
	File *os.File
}

func (fs *TeeFileStreamer) Run() {
	var s string
	for {
		s = <- fs.In
		switch (s) {
		case ZERO_STR:
			fs.Out <- s
			return
		default:
			fs.File.WriteString(s)
			fs.Out <- s
		}
	}
}

type TokenValueFileStreamer struct {
	In   chan(lexer.TokenValue)
	Out  chan(lexer.TokenValue)
	File *os.File
}

func (fs *TokenValueFileStreamer) Run() {
	var s lexer.TokenValue
	for {
		s = <- fs.In
		switch {
		case s.Token == lexer.EndToken:
			return
		default:
			fs.File.WriteString(s.Value)
			fs.Out <- s
		}
	}
}
