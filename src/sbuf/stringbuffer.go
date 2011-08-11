// Simple string buffer
package sbuf

import "strings"

const initialSize = 8

type StringBuffer []string

func (p *StringBuffer) realloc(length, capacity int) (b []string) {
	if capacity < initialSize {
		capacity = initialSize
	}
	if capacity < length {
		capacity = length
	}
	b = make(StringBuffer, length, capacity)
	copy(b, *p)
	*p = b
	return
}

func (p *StringBuffer) expand(i, n int) {
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

func (p *StringBuffer) Reset() {
	a := *p
	n := len(a)

	var zero string
	for k := 0; k < n; k++ {
		a[k] = zero
	}

	*p = a[0:0]
}

func (p *StringBuffer) Join(sep string) (string) {
	return strings.Join(*p, sep)
}

func (p *StringBuffer) Empty() (bool) {
	return p.Len() == 0
}

func (p *StringBuffer) At(index int) (string) {
	return (*p)[index]
}

func (p *StringBuffer) Len() (int) {
	return len(*p)
}

func (p *StringBuffer) Set(index int, str string) {
	(*p)[index] = str
}

func (p *StringBuffer) Insert(index int, str string) {
	p.expand(index, 1)
	(*p)[index] = str
}

func (p *StringBuffer) Push(s string) {
	p.Insert(len(*p), s)
}

func (p *StringBuffer) Delete(i int) {
	a := *p
	n := len(a)

	copy(a[i:n-1], a[i+1:n])
	var zero string
	a[n-1] = zero // support GC, zero out entry
	*p = a[0:n-1]
}
