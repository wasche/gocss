/**
 * @author wasche
 * @since 2011.08.28
 */
package main

import (
	"os"
	"bufio"
	"fmt"
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

type OutputFileStreamer struct {
	In  chan(string)
	Out *os.File
	Eof chan(int)
}

var ZERO_STR string

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
