package lexer

// Token types
type Token int
const (
	_ = iota
	Whitespace Token = iota
	Comment
	String
	Identifier
	Number
	Op
	// operator types
	Hash
	Percent
	LeftParen
	RightParen
	Match
	Child
	At
	Colon
	Semicolon
	LeftBrace
	RightBrace
	LeftBracket
	RightBracket
	Plus
	Star
	Period
	Comma
	Tildae
	Pipe
	Backslash
	Dollar
	Carrot
	Bang
)

func (b Token) String() string {
	switch b {
	case Whitespace:
		return "Whitespace"
	case Comment:
		return "Comment"
	case String:
		return "String"
	case Identifier:
		return "Identifier"
	case Number:
		return "Number"
	case Hash:
		return "#"
	case Percent:
		return "%"
	case LeftParen:
		return "("
	case RightParen:
		return ")"
	case Match:
		return "="
	case Child:
		return ">"
	case At:
		return "@"
	case Colon:
		return ":"
	case Semicolon:
		return ";"
	case LeftBrace:
		return "{"
	case RightBrace:
		return "}"
	case LeftBracket:
		return "["
	case RightBracket:
		return "]"
	case Plus:
		return "+"
	case Star:
		return "*"
	case Period:
		return "."
	case Comma:
		return ","
	case Tildae:
		return "~"
	case Pipe:
		return "|"
	case Backslash:
		return "\\"
	case Dollar:
		return "$"
	case Carrot:
		return "^"
	case Bang:
		return "!"
	}
	return "Unknown token"
}

var TokenMap = map[int] Token {
	'{': LeftBrace,
	'}': RightBrace,
	'[': LeftBracket,
	']': RightBracket,
	'(': LeftParen,
	')': RightParen,
	'+': Plus,
	'*': Star,
	'=': Match,
	',': Comma,
	';': Semicolon,
	':': Colon,
	'>': Child,
	'~': Tildae,
	'|': Pipe,
	'\\': Backslash,
	'%': Percent,
	'$': Dollar,
	'#': Hash,
	'@': At,
	'^': Carrot,
	'!': Bang,
	'/': Comment,
	'_': Identifier,
	'a': Identifier,
	'b': Identifier,
	'c': Identifier,
	'd': Identifier,
	'e': Identifier,
	'f': Identifier,
	'g': Identifier,
	'h': Identifier,
	'i': Identifier,
	'j': Identifier,
	'k': Identifier,
	'l': Identifier,
	'm': Identifier,
	'n': Identifier,
	'o': Identifier,
	'p': Identifier,
	'q': Identifier,
	'r': Identifier,
	's': Identifier,
	't': Identifier,
	'u': Identifier,
	'v': Identifier,
	'w': Identifier,
	'x': Identifier,
	'y': Identifier,
	'z': Identifier,
	'A': Identifier,
	'B': Identifier,
	'C': Identifier,
	'D': Identifier,
	'E': Identifier,
	'F': Identifier,
	'G': Identifier,
	'H': Identifier,
	'I': Identifier,
	'J': Identifier,
	'K': Identifier,
	'L': Identifier,
	'M': Identifier,
	'N': Identifier,
	'O': Identifier,
	'P': Identifier,
	'Q': Identifier,
	'R': Identifier,
	'S': Identifier,
	'T': Identifier,
	'U': Identifier,
	'V': Identifier,
	'W': Identifier,
	'X': Identifier,
	'Y': Identifier,
	'Z': Identifier,
	'-': Number,
	'.': Number,
	'0': Number,
	'1': Number,
	'2': Number,
	'3': Number,
	'4': Number,
	'5': Number,
	'6': Number,
	'7': Number,
	'8': Number,
	'9': Number,
	'"': String,
	'\'': String,
	' ': Whitespace,
	'\t': Whitespace,
	'\n': Whitespace,
}

