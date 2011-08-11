// Minification parser
package parser

import (
	"./lexer"
	"./sbuf"
	"os"
	"strings"
)

const MS_ALPHA = "progid:dximagetransform.microsoft.alpha(opacity="

var (
	UNITS           = map[string] bool {
		"px": true,
		"em": true,
		"pt": true,
		"in": true,
		"cm": true,
		"mm": true,
		"pc": true,
		"ex": true,
		"%":  true,
	}
	KEYWORDS        = map[string] bool {
		"normal":     true,
		"bold":       true,
		"italic":     true,
		"serif":      true,
		"sans-serif": true,
		"fixed":      true,
	}
	BOUNDARY_OPS    = map[lexer.Token] bool {
		lexer.LeftBrace:  true,
		lexer.RightBrace: true,
		lexer.Child:      true,
		lexer.Semicolon:  true,
		lexer.Colon:      true,
		lexer.Comma:      true,
		lexer.Comment:    true,
	}
	NONE_PROPERTIES = make(map[string]bool, 15)
)

func init() {
	NONE_PROPERTIES["outline"] = true
	props := [...]string{"border", "margin", "padding"}
	edges := [...]string{"top", "left", "right", "bottom"}
	for _, p := range props {
		for _, e := range edges {
			NONE_PROPERTIES[p + "-" + e] = true
		}
	}
}

type Parser struct {
	lastToken   lexer.Token
	lastValue   string
	property    string
	ruleBuffer  sbuf.StringBuffer
	valueBuffer sbuf.StringBuffer
	rgbBuffer   sbuf.StringBuffer
	pending     string
	inRule      bool
	space       bool
	charset     bool
	at          bool
	ie5mac      bool
	rgb         bool
	checkSpace  int
}

func (p *Parser) output(str string) {
	os.Stdout.WriteString(str)
}

func (p *Parser) dump(str string) {
	p.ruleBuffer.Push(p.pending)
	p.ruleBuffer.Push(str)
	p.output(p.ruleBuffer.Join(""))
	p.ruleBuffer.Reset()
	p.pending = ""
}

func (p *Parser) write(str string) {
	if len(str) == 0 { return }
	if len(str) >= 3 && str[0:3] == "/*!" && p.ruleBuffer.Empty() {
		p.output(str)
		return
	}
	p.ruleBuffer.Push(str)
	if str == "}" {
		// check for empty rule
		s := p.ruleBuffer.Join("")
		if s[len(s)-2:] != "{}" {
			p.output(s)
		}
		p.ruleBuffer.Reset()
	}
}

func (p *Parser) buffer(str string) {
	if len(p.pending) > 0 {
		p.write(p.pending)
	}
	p.pending = str
}

func (p *Parser) q(str string) {
	switch {
	case p.property == "":
		p.buffer(str)
	default:
		p.valueBuffer.Push(str)
	}
}

func (p *Parser) collapseZeroes() {
	t := p.valueBuffer.Join("")
	p.valueBuffer.Reset()
	_, isNone := NONE_PROPERTIES[p.property]
	switch {
	case t == "0 0" || t == "0 0 0" || t == "0 0 0 0":
		p.buffer("0")
		if p.property == "background-positon" || p.property == "-webkit-transform-origin" || p.property == "-moz-transform-origin" {
			p.buffer(" 0")
		}
	case t == "none" && (p.property == "background" || isNone):
		p.buffer("0")
	default:
		p.buffer(t)
	}
}

func (p *Parser) Token(token lexer.Token, value string) {
	//os.Stderr.WriteString("token: "+token+", value: "+value+"\n")

	if p.rgb {
		switch token {
		case lexer.Number:
			// TODO
		case lexer.LeftParen:
			if p.lastToken == lexer.Number { p.q(" ") }
			p.q("#")
			p.rgbBuffer.Reset()
		case lexer.RightParen:
			// TODO
		}
		return
	}

	if token == lexer.Whitespace {
		p.space = true
		return
	}

	if token == lexer.Comment {
		// comments are only needed in a few places:
		switch {
		// 1) special comments /*! ... */
		case value[2:3] == "!":
			p.q(value)
			p.lastToken = token
			p.lastValue = value
		case value[len(value)-3:len(value)-2] == "\\":
			p.q("/*\\*/")
			p.lastToken = token
			p.lastValue = value
			p.ie5mac = true
		case p.ie5mac:
			p.q("/**/")
			p.lastToken = token
			p.lastValue = value
			p.ie5mac = false
		case p.lastToken == lexer.Child:
			p.q("/**/")
			p.lastToken = token
			p.lastValue = value
		}
		return
	}

	// most whitespace isn't needed, but make sure we have space between values
	// for multivalue properties
	// margin: 5px 5px;
	switch {
	case p.lastToken == lexer.Number && (token == lexer.Hash || token == lexer.Number):
		p.q(" ")
		p.space = false
	case (token == lexer.Number || token == lexer.Identifier || token == lexer.Hash) &&
			(p.lastToken == lexer.Identifier || p.lastToken == lexer.Percent || p.lastToken == lexer.RightParen):
		p.q(" ")
		p.space = false
	case p.inRule && token == lexer.Identifier && p.lastToken == lexer.RightParen:
		p.q(" ")
		p.space = false
	}

	// rgb()
	if token == lexer.Identifier && value == "rgb" {
		p.space = false
		p.rgb = true
		return
	}

	switch {
	case token == lexer.At:
		p.q(value)
		p.at = true
	case p.inRule && token == lexer.Colon && len(p.property) == 0:
		p.q(value)
		if len(p.lastValue) == 0 { p.property = strings.ToLower(p.lastValue) }
		p.valueBuffer.Reset()
	// first-letter and first-line must be followed by a space
	case !p.inRule && p.lastToken == lexer.Colon && (value == "first-letter" || value == "first-line"):
		p.q(value)
		p.q(" ")
	case token == lexer.Semicolon:
		switch {
		case p.at:
			p.at = false
			switch {
			default:
				p.dump(value)
			case p.ruleBuffer.At(1) == "charset":
				switch {
				case p.charset:
					p.ruleBuffer.Reset()
					p.pending = ""
				default:
					p.charset = true
					p.dump(value)
				}
			}
		case p.lastToken == lexer.Semicolon:
			// skip
			return
		default:
			p.collapseZeroes()
			p.valueBuffer.Reset()
			p.property = ""
			p.q(value)
		}
	case token == lexer.LeftBrace:
		// TODO
	case token == lexer.RightBrace:
		// TODO
	case !p.inRule:
		// TODO
	case token == lexer.Number && len(value) > 2 && value[:2] == "0.":
		p.q(value[2:])
	case token == lexer.String && p.property == "-ms-filter":
		if strings.ToLower(value[1:len(MS_ALPHA)+1]) == MS_ALPHA {
			c := value[0:1]
			a := value[len(MS_ALPHA)+1:len(value)-2]
			p.q(c)
			p.q("alpha(opacity=")
			p.q(a)
			p.q(")")
			p.q(c)
		} else {
			p.q(value)
		}
	case token == lexer.Match:
		p.q(value)
		if strings.ToLower(p.valueBuffer.Join("")) == MS_ALPHA {
			p.buffer("alpha(opacity=")
			p.valueBuffer.Reset()
		}
	default:
		t := strings.ToLower(value)
		switch {
		// values of 0 don't need a unit
		case p.lastToken == lexer.Number && p.lastValue == "0" &&
				(token == lexer.Percent || token == lexer.Identifier):
			if UNITS[value] {
				p.q(" ")
				p.q(value)
			}
		// use 0 instead of none
		case value == "none" && p.lastToken == lexer.Colon && NONE_PROPERTIES[p.property]:
			p.q("0")
		// force properties to lower case for better gzip compression
		case token == lexer.Identifier && p.lastToken == lexer.Colon:
			switch {
			// #aabbcc
			case p.lastToken == lexer.Hash:
				if len(value) == 6 &&
						t[0] == t[1] &&
						t[2] == t[3] &&
						t[4] == t[5] {
					p.q(t[1:3])
					p.q(t[4:5])
				} else {
					p.q(t)
				}
			case len(p.property) == 0 || KEYWORDS[t]:
				p.q(t)
			default:
				p.q(value)
			}
		}
	}

	p.lastToken = token
	p.lastValue = value
	p.space = false
}

func (p *Parser) End() {
	p.write(p.pending)
	if !p.ruleBuffer.Empty() {
		p.output(p.ruleBuffer.Join(""))
	}
}

