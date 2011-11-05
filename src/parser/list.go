/**
 * A simple list of TokenValues.
 * @author wasche
 * @since 2011.08.30
 */
package parser

import ("./lexer")

const initialSize = 8
var BLANK lexer.TokenValue

type TokenValueList []lexer.TokenValue

func (p *TokenValueList) realloc(length, capacity int) (b []lexer.TokenValue) {
	if capacity < initialSize {
		capacity = initialSize
	}
	if capacity < length {
		capacity = length
	}
	b = make(TokenValueList, length, capacity)
	copy(b, *p)
	*p = b
	return
}

func (p *TokenValueList) expand(i, n int) {
	a := *p

	len0 := len(a)
	len1 := len0 + n
	if len1 <= cap(a) {
		// enough space - just expand
		a = a[0:len1]
	} else {
		// not enough space - double capacity
		capb := cap(a) * 2
		if capb < len1 {
			// still not enough - use required length
			capb = len1
		}
		a = p.realloc(len1, capb)
	}

	// make a hole
	for j := len0 - 1; j >= i; j-- {
		a[j+n] = a[j]
	}

	*p = a
}

func (p *TokenValueList) Reset() {
	a := *p
	n := len(a)

	for k := 0; k < n; k++ {
		a[k] = BLANK
	}

	*p = a[0:0]
}

func (p *TokenValueList) Join(sep string) (string) {
	if len(p) == 0 {
		return ""
	}
	if len(p) == 1 {
		return p[0].Value
	}
	n := len(sep) * (len(p) - 1)
	for i := 0; i < len(p); i++ {
		n += len(a[i].Value)
	}

	b := make([]byte, n)
	bp := copy(b, a[0].Value)
	for _, s := range a[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], s.Value)
	}
	return string(b)
}

func (p *TokenValueList) Empty() (bool) {
	return p.Len() == 0
}

func (p *TokenValueList) At(index int) (lexer.TokenValue) {
	return (*p)[index]
}

func (p *TokenValueList) Len() (int) {
	return len(*p)
}

func (p *TokenValueList) Set(index int, tv lexer.TokenValue) {
	(*p)[index] = tv
}

func (p *TokenValueList) Insert(index int, tv lexer.TokenValue) {
	p.expand(index, 1)
	(*p)[index] = tv
}

func (p *TokenValueList) Push(tv lexer.TokenValue) {
	p.Insert(len(*p), tv)
}

func (p *TokenValueList) PushAll(list TokenValueList) {
	// TODO change to use copy for better performance
	for i := 0; i < len(list); i++ {
		p.Push(list.At[i])
	}
}

func (p *TokenValueList) Delete(i int) {
	a := *p
	n := len(a)

	copy(a[i:n-1], a[i+1:n])
	a[n-1] = BLANK // support GC, zero out entry
	*p = a[0:n-1]
}
