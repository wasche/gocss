// Minification parser
package parser

import (
	"./lexer"
	//"fmt"
	"os"
	"bytes"
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
		"%": true,
	}
	KEYWORDS        = map[string] bool {
		"normal": true,
		"bold": true,
		"italic": true,
		"serif": true,
		"sans-serif": true,
		"fixed": true,
	}
	BOUNDARY_OPS    = map[lexer.Token] bool {
		lexer.LeftBrace: true,
		lexer.RightBrace: true,
		lexer.Child: true,
		lexer.Semicolon: true,
		lexer.Colon: true,
		lexer.Comma: true,
		lexer.Comment: true,
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
	ruleBuffer  bytes.Buffer
	valueBuffer bytes.Buffer
	rgbBuffer   bytes.Buffer
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
	p.ruleBuffer.WriteString(p.pending)
	p.ruleBuffer.WriteString(str)
	p.output(p.ruleBuffer.String())
	p.ruleBuffer.Reset()
	p.pending = ""
}

func (p *Parser) write(str string) {
	if len(str) == 0 { return }
	if len(str) >= 3 && str[0:3] == "/*!" && p.ruleBuffer.Len() == 0 {
		p.output(str)
		return
	}
	p.ruleBuffer.WriteString(str)
	if str == "}" {
		// check for empty rule
		s := p.ruleBuffer.String()
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
		p.valueBuffer.WriteString(str)
	}
}

func (p *Parser) collapseZeroes() {
	t := p.valueBuffer.String()
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
	//fmt.Fprintf(os.Stderr, "token: %s, value: %s\n", token, value)

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
		// TODO
	case !p.inRule && p.lastToken == lexer.Colon && (value == "first-letter" || value == "first-line"):
		p.q(value)
		p.q(" ")
	case token == lexer.Semicolon:
		// TODO
	case token == lexer.LeftBrace:
		// TODO
	case token == lexer.RightBrace:
		// TODO
	case !p.inRule:
		// TODO
	case token == lexer.Number && len(value) > 2 && value[:2] == "0.":
		p.q(value[2:])
	case token == lexer.String && p.property == "-ms-filter":
		// TODO
	case token == lexer.Match:
		// TODO
	default:
		// TODO
	}

	p.lastToken = token
	p.lastValue = value
	p.space = false
}

func (p *Parser) End() {
	p.write(p.pending)
	if p.ruleBuffer.Len() > 0 {
		p.output(p.ruleBuffer.String())
	}
}

