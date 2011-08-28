// Simple streaming lexer for CSS
package lexer

import "bytes"

type HandlerFn func(*Lexer, int)

type TokenValue struct {
	Token Token
	Value string
}

type Lexer struct {
	In chan(int)
	Out chan(TokenValue)
	token bytes.Buffer
	prev int
	lastToken string
	handler HandlerFn
}

func isDigit(c int) bool {
	return c >= '0' && c <= '9'
}

func isNameChar(c int) bool {
	return c == '_' || c == '-' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
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
		lex.tokenize(c)
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
	switch {
	case lex.token.Len() == 0 || lex.prev == '\\':
		lex.token.WriteRune(c)
	case c != f:
		lex.token.WriteRune(c)
	default:
		lex.token.WriteRune(c)
		lex.next(String)
		lex.token.Reset()
		lex.handler = nil
	}
}

func (lex *Lexer) identifier(c int) {
	switch {
	case isNameChar(c) || isDigit(c):
		lex.token.WriteRune(c)
	default:
		lex.next(Identifier)
		lex.token.Reset()
		lex.handler = nil
		lex.tokenize(c)
	}
}

func (lex *Lexer) number(c int) {
	nondigit := !isDigit(c)
	point := '.' == lex.prev
	switch {
	// .2em or .classname ?
	case point && nondigit:
		lex.next(Period)
		lex.token.Reset()
		lex.handler = nil
		lex.tokenize(c)
	// -2px or -moz-something
	case '-' == lex.prev && nondigit && c != '.':
		lex.handler = (*Lexer).identifier
		lex.handler(lex, c)
	case !nondigit || (!point && ('.' == c || '-' == c)):
		lex.token.WriteRune(c)
	case lex.lastToken == "#":
		lex.handler = (*Lexer).identifier
		lex.handler(lex, c)
	default:
		lex.next(Number)
		lex.token.Reset()
		lex.handler = nil
		lex.tokenize(c)
	}
}

func (lex *Lexer) operator(c int) {
	switch {
	case '=' == c:
		lex.token.WriteRune(c)
		lex.next(Match)
		lex.token.Reset()
		lex.handler = nil
	case lex.token.Len() == 0:
		lex.token.WriteRune(c)
	default:
		lex.next(TokenMap[lex.prev])
		lex.token.Reset()
		lex.handler = nil
		lex.tokenize(c)
	}
}

func (lex *Lexer) next(t Token) {
	value := lex.token.String()
	lex.Out <- TokenValue{t, value}
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

func (lex *Lexer) tokenize(c int) {
	if lex.handler == nil {
		token := TokenMap[c]
		lex.handler = lex.handlerForToken(token)
	}
	lex.handler(lex, c)
	lex.prev = c
}

func (lex *Lexer) end() {
	switch {
	// finish the current token
	case lex.handler != nil:
		lex.handler(lex, -1)
	// still something in buffer, assuming whitespace
	case lex.token.Len() > 0:
		lex.next(Whitespace)
	}
	lex.Out <- TokenValue{EndToken, ""}
}

func (lex *Lexer) Run() {
	var c int
	for {
		c = <- lex.In
		if c == -1 {
			lex.end()
			return
		} else {
			lex.tokenize(c)
		}
	}
}
