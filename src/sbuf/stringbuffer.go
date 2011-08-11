// Simple string buffer
package sbuf

import "strings"

const initialSize = 8

type StringBuffer struct {
	buffer []string
	length int
}

func (p *StringBuffer) Push(s string) {
	switch {
	case p.length == cap(p.buffer):
		p.buffer = append(p.buffer, s)
	default:
		p.buffer[p.length] = s
		p.length++
	}
}

func (p *StringBuffer) Reset() {
	p.length = 0
}

func (p *StringBuffer) Join(sep string) (string) {
	return strings.Join(p.buffer, sep)
}

func (p *StringBuffer) Empty() (bool) {
	return p.length == 0
}

func (p *StringBuffer) At(index int) (string) {
	return p.buffer[index]
}
