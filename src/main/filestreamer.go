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

func CreateInputFileStreamer(file *os.File) (ifs *InputFileStreamer, out chan(int)) {
	out = make(chan(int))
	ifs = &InputFileStreamer{In: file, Out: out}
	return
}

var ZERO_STR string

type OutputFileStreamer struct {
	In  chan(lexer.TokenValue)
	Out *os.File
	Eof chan(int)
}

func (fs *OutputFileStreamer) Run() {
	var s lexer.TokenValue
	for {
		s = <- fs.In
		switch {
		case lexer.EndToken == s.Token:
			if fs.Eof != nil { fs.Eof <- 0 }
			return
		default:
			fs.Out.WriteString(s.Value)
		}
	}
}

func CreateOutputFileStreamer (in chan(lexer.TokenValue), file *os.File) (ofs *OutputFileStreamer, eof chan(int)) {
	eof = make(chan(int))
	ofs = &OutputFileStreamer{In: in, Out: file, Eof: eof}
	return
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
		fs.Out <- s
		switch (s) {
		case ZERO_STR:
			return
		default:
			fs.File.WriteString(s)
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
		fs.Out <- s
		switch {
		case s.Token == lexer.EndToken:
			return
		default:
			fs.File.WriteString(s.Value)
		}
	}
}

func CreateTokenValueFileStreamer(in chan(lexer.TokenValue), file *os.File) (tvfs *TokenValueFileStreamer, out chan(lexer.TokenValue)) {
	out = make(chan(lexer.TokenValue))
	tvfs = &TokenValueFileStreamer{In: in, File: file, Out: out}
	return
}

type ChannelSplitter struct {
  In   chan(lexer.TokenValue)
  Out1 chan(lexer.TokenValue)
  Out2 chan(lexer.TokenValue)
}

func (fs *ChannelSplitter) Run() {
	var s lexer.TokenValue
	for {
		s = <- fs.In
		fs.Out1 <- s
		fs.Out2 <- s
		if s.Token == lexer.EndToken { return }
	}
}

func CreateChannelSplitter(in chan(lexer.TokenValue)) (cs *ChannelSplitter, out1 chan(lexer.TokenValue), out2 chan(lexer.TokenValue)) {
	out1 = make(chan(lexer.TokenValue))
	out2 = make(chan(lexer.TokenValue))
	cs = &ChannelSplitter{In: in, Out1: out1, Out2: out2}
	return
}
