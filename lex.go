package plurals

import "fmt"

func Lex(s string) (tokens []Token, err error) {
	var (
		pos   = 0
		siz   = len(s)
		token Token
	)
	for {
		token, pos = readToken(s, pos, siz)
		switch token.Type {
		case TokenTypeEOF:
			return
		case TokenTypeERR:
			err = fmt.Errorf("error read token: %s, input: %q", token.Value, s)
			return
		}
		tokens = append(tokens, token)
	}
}

func readToken(s string, pos, siz int) (token Token, newPos int) {
	// 跳过空白字符
	var ch byte
	for {
		if pos >= siz {
			return Token{Type: TokenTypeEOF}, siz
		}
		ch = s[pos]
		pos++
		if ch != ' ' && ch != '\t' {
			break
		}
	}
	val := []byte{ch}
	switch ch {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		start := pos - 1
		num := int64(ch - '0')
		for pos < siz && s[pos] >= '0' && s[pos] <= '9' {
			ch = s[pos]
			pos++
			val = append(val, ch)
			num *= 10
			num += int64(ch - '0')
		}
		return Token{
			Type:   TokenTypeNUM,
			Value:  string(val),
			Number: num,
			Start:  start,
			End:    pos,
		}, pos
	case '!':
		if pos < siz && s[pos] == '=' {
			pos++
			return Token{
				Type:  TokenTypeEQU,
				Value: "!=",
				Start: pos - 2,
				End:   pos,
			}, pos
		}
		return Token{
			Type:  TokenTypeLGC,
			Value: "!",
			Start: pos - 1,
			End:   pos,
		}, pos
	case '&', '|', '=':
		if pos < siz && s[pos] == ch {
			pos++
			val = append(val, ch)
			return Token{
				Type:  ch2Typ[ch],
				Value: string(val),
				Start: pos - 2,
				End:   pos,
			}, pos
		}
		return Token{
			Type: TokenTypeERR,
			Value: fmt.Sprintf("at column [%d:%d]: expected '%c%c'",
				pos-1, pos, ch, ch),
		}, pos
	case '<', '>':
		if pos < siz && s[pos] == '=' {
			ch = s[pos]
			pos++
			val = append(val, ch)
			return Token{
				Type:  TokenTypeCMP,
				Value: string(val),
				Start: pos - 2,
				End:   pos,
			}, pos
		}
		return Token{
			Type:  TokenTypeCMP,
			Value: string(val),
			Start: pos - 1,
			End:   pos,
		}, pos
	case '*', '/', '%', '+', '-', '?', ':', 'n', '(', ')', ';':
		return Token{
			Type:  ch2Typ[ch],
			Value: string(val),
			Start: pos - 1,
			End:   pos,
		}, pos
	default:
		return Token{
			Type: TokenTypeERR,
			Value: fmt.Sprintf("at column [%d:%d] unexpected '%c'",
				pos-1, pos, ch),
		}, pos
	}
}

var ch2Typ = map[byte]TokenType{
	'=': TokenTypeEQU,
	'|': TokenTypeLGC,
	'&': TokenTypeLGC,
	'+': TokenTypeADD,
	'-': TokenTypeADD,
	'*': TokenTypeMUL,
	'/': TokenTypeMUL,
	'%': TokenTypeMUL,
	':': TokenTypeCOL,
	'n': TokenTypeIDN,
	'?': TokenTypeQST,
	'(': TokenTypeLPA,
	')': TokenTypeRPA,
	';': TokenTypeCOM,
}
