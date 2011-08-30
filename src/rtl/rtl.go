/**
 * @author wasche
 * @since 2011.08.29
 */
package rtl

import (
	"./lexer"
	"os"
)

type Converter struct {
	In  chan(lexer.TokenValue)
	Out *os.File
}

func CreateConverter(in chan(lexer.TokenValue), file *os.File) (c *Converter) {
	c = &Converter{In: in, Out: file}
	return
}

func (c *Converter) Run() {
	var tv lexer.TokenValue
	for {
		tv = <- c.In
		if tv.Token == lexer.EndToken {
			return
		} else {
			c.Out.WriteString(tv.Value)
		}
	}
}
