// Simple streaming lexer for CSS
package lexer

import "bytes"

type HandlerFn func(*Lexer, int)

type Parser interface {
	Token(Token, string)
	End()
}

type Lexer struct {
	token bytes.Buffer
	prev int
	lastToken string
	handler HandlerFn
	parser Parser
}

func isDigit(c int) bool {
	return c >= '0' && c <= '9'
}

func isNameChar(c int) bool {
	return c == '_' || c == '-' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func (lex *Lexer) SetParser(p Parser) {
	lex.parser = p
}

func (lex *Lexer) whitespace(c int) {
	switch {
	case c == ' ' || c == '\t' || c == '\n':
		lex.token.WriteRune(c)
		return
	default:
		lex.next(Whitespace)
		lex.token.Reset()
		lex.handler = nil
		lex.Tokenize(c)
	}
}

func (lex *Lexer) comment(c int) {
	switch {
	case lex.token.Len() == 1 && c != '*':
		lex.next(Op)
		lex.token.Reset()
		lex.token.WriteRune(c)
		lex.handler = nil
		return
	case c == '/' && lex.prev == '*':
		lex.token.WriteRune(c)
		lex.next(Comment)
		lex.token.Reset()
		lex.handler = nil
	default:
		lex.token.WriteRune(c)
	}
}

func (lex *Lexer) str(c int) {
	f, _, _ := lex.token.ReadRune()
	lex.token.UnreadRune()
	if lex.token.Len() == 0 || lex.prev == '\\' {
		lex.token.WriteRune(c)
	} else if c == f {
		lex.token.WriteRune(c)
	} else {
		lex.token.WriteRune(c)
		lex.next(String)
		lex.token.Reset()
		lex.handler = nil
	}
}

func (lex *Lexer) identifier(c int) {
	if isNameChar(c) || isDigit(c) {
		lex.token.WriteRune(c)
	} else {
		lex.next(Identifier)
		lex.token.Reset()
		lex.handler = nil
		lex.Tokenize(c)
	}
}

func (lex *Lexer) number(c int) {
	nondigit := !isDigit(c)
	point := '.' == lex.prev
	// .2em or .classname ?
	if point && nondigit {
		lex.next(Period)
		lex.token.Reset()
		lex.handler = nil
		lex.Tokenize(c)
	// -2px or -moz-something
	} else if '-' == lex.prev && nondigit {
		lex.handler = (*Lexer).identifier
		lex.handler(lex, c)
	} else if !nondigit || (!point && ('.' == c || '-' == c)) {
		lex.token.WriteRune(c)
	} else if lex.lastToken == "#" {
		lex.handler = (*Lexer).identifier
		lex.handler(lex, c)
	} else {
		lex.next(Number)
		lex.token.Reset()
		lex.handler = nil
		lex.Tokenize(c)
	}
}

func (lex *Lexer) operator(c int) {
	if '=' == c {
		lex.token.WriteRune(c)
		lex.next(Match)
		lex.token.Reset()
		lex.handler = nil
	} else if lex.token.Len() == 0 {
		lex.token.WriteRune(c)
	} else {
		lex.next(TokenMap[c])
		lex.token.Reset()
		lex.handler = nil
		lex.Tokenize(c)
	}
}

func (lex *Lexer) next(t Token) {
	value := lex.token.String()
	lex.parser.Token(t, value)
	lex.lastToken = value
}

func (lex *Lexer) handlerForToken(t Token) (fn HandlerFn) {
	switch t {
	case Whitespace:
		return (*Lexer).whitespace
	case Comment:
		return (*Lexer).comment
	case String:
		return (*Lexer).str
	case Identifier:
		return (*Lexer).identifier
	case Number:
		return (*Lexer).number
	}
	return (*Lexer).operator
}

func (lex *Lexer) Tokenize(c int) {
	if lex.handler == nil {
		token := TokenMap[c]
		lex.handler = lex.handlerForToken(token)
	}
	lex.handler(lex, c)
	lex.prev = c
}

func (lex *Lexer) End() {
	if lex.handler != nil {
		lex.handler(lex, -1)
	} else if lex.token.Len() > 0 {
		value := lex.token.String()
		lex.parser.Token(Whitespace, value)
	}
	lex.parser.End()
}
