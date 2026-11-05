package db0302

import (
	"errors"
	"fmt"
	"strings"
)

type Parser struct {
	buf string
	pos int
}

func NewParser(s string) Parser {
	return Parser{buf: s, pos: 0}
}

func isSpace(ch byte) bool {
	switch ch {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	}
	return false
}
func isAlpha(ch byte) bool {
	return 'a' <= (ch|32) && (ch|32) <= 'z'
}
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
func isNameStart(ch byte) bool {
	return isAlpha(ch) || ch == '_'
}
func isNameContinue(ch byte) bool {
	return isAlpha(ch) || isDigit(ch) || ch == '_'
}
func isSeparator(ch byte) bool {
	return ch < 128 && !isNameContinue(ch)
}

func (p *Parser) skipSpaces() {
	for p.pos < len(p.buf) && isSpace(p.buf[p.pos]) {
		p.pos += 1
	}
}

func (p *Parser) tryKeyword(kw string) bool {
	p.skipSpaces()
	if !(p.pos+len(kw) <= len(p.buf) && strings.EqualFold(p.buf[p.pos:p.pos+len(kw)], kw)) {
		return false
	}
	if p.pos+len(kw) < len(p.buf) && !isSeparator(p.buf[p.pos+len(kw)]) {
		return false
	}
	p.pos += len(kw)
	return true
}

func (p *Parser) tryName() (string, bool) {
	p.skipSpaces()
	start, cur := p.pos, p.pos
	if !(cur < len(p.buf) && isNameStart(p.buf[cur])) {
		return "", false
	}
	cur++
	for cur < len(p.buf) && isNameContinue(p.buf[cur]) {
		cur++
	}
	p.pos = cur
	return p.buf[start:cur], true
}

func (p *Parser) parseValue(out *Cell) error {
	p.skipSpaces()
	if p.pos >= len(p.buf) {
		return errors.New("expect value")
	}
	ch := p.buf[p.pos]
	if ch == '"' || ch == '\'' {
		return p.parseString(out)
	} else if isDigit(ch) || ch == '-' || ch == '+' {
		return p.parseInt(out)
	} else {
		return errors.New("expect value")
	}
}

func (p *Parser) parseString(out *Cell) error {
	startQoutePos := p.pos
	p.pos += 1 //points to first char
	s := ""
	for p.pos < len(p.buf) {
		if p.buf[p.pos] == p.buf[startQoutePos] { //reached the end
			p.pos += 1
			break
		}
		if p.buf[p.pos] == '\\' { // escaped characters
			p.pos += 1
			switch p.buf[p.pos] {
			case '\'':
				s += "'"
				p.pos += 1
				continue
			case '"':
				s += "\""
				p.pos += 1
				continue
			case '\\':
				s += "\\"
				p.pos += 1
				continue
			}
		}
		s += string(p.buf[p.pos])
		p.pos += 1
	}
	out.Type = TypeStr
	out.Str = []byte(s)

	return nil
}

func (p *Parser) parseInt(out *Cell) (err error) {
	start := p.pos
	if !isDigit(p.buf[p.pos]) {
		p.pos += 1
		start = p.pos
	}

	for isDigit(p.buf[p.pos]) {
		p.pos += 1
	}

	var n int64
	for idx := range p.buf[start:p.pos] {
		fmt.Println("idx: ", idx, int64(p.buf[start+idx]-'0'))
		n = n*10 + int64(p.buf[start+idx]-'0')
	}

	out.Type = TypeI64
	if p.buf[start-1] == '-' {
		out.I64 = -n
	} else {
		out.I64 = n
	}
	return nil
}

func (p *Parser) isEnd() bool {
	p.skipSpaces()
	return p.pos >= len(p.buf)
}

// QzBQWVJJOUhU https://trialofcode.org/
