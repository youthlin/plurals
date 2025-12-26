package plurals

import "fmt"

func Compile(s string) (Expression, error) {
	tokens, err := Lex(s)
	if err != nil {
		return nil, err
	}
	return parse(tokens)
}

func parse(tokens []Token) (node Expression, err error) {
	index := 0
	total := len(tokens)
	if total == 0 {
		err = fmt.Errorf("empty token")
		return
	}
	index, node, err = parseExpression(tokens, total, index)
	if err != nil {
		return
	}
	for _, ok := get(tokens, total, index); ok; {
		_, index, err = consume(tokens, total, index, TokenTypeCOM, ";")
		if err != nil {
			return
		}
	}
	return node, err
}

func consume(tokens []Token, total, idx int, expType TokenType, expVal string) (
	token Token,
	index int,
	err error,
) {
	index = idx
	if index >= total {
		err = fmt.Errorf("expected `%v(%v)` after token %v",
			expVal, expType, tokens[total-1])
		return
	}
	token = tokens[index]
	if token.Type != expType {
		err = fmt.Errorf("expected `%v(%v)`, but got %v",
			expVal, expType, token)
		return
	}
	if expVal != "" && token.Value != expVal {
		err = fmt.Errorf("expected `%v(%v)`, but got %v",
			expVal, expType, token)
		return
	}
	index = index + 1
	return
}

func get(tokens []Token, total, index int) (Token, bool) {
	if index < total {
		return tokens[index], true
	}
	return Token{}, false
}

func parseExpression(tokens []Token, total, idx int) (index int, node Expression, err error) {
	return parseTernary(tokens, total, idx)
}

func parseTernary(tokens []Token, total, idx int) (index int, node Expression, err error) {
	index = idx
	// logicOr ( '?' exp ':' exp )?
	index, node, err = parseLogicOr(tokens, total, index)
	if err != nil {
		return
	}
	if token, ok := get(tokens, total, index); ok && token.Type == TokenTypeQST {
		_, index, err = consume(tokens, total, index, TokenTypeQST, "?")
		if err != nil {
			return
		}
		var branchTrue Expression
		index, branchTrue, err = parseExpression(tokens, total, index)
		if err != nil {
			return
		}
		_, index, err = consume(tokens, total, index, TokenTypeCOL, ":")
		if err != nil {
			return
		}
		var branchFalse Expression
		index, branchFalse, err = parseExpression(tokens, total, index)
		if err != nil {
			return
		}
		node = &TernaryNode{
			Condition:   node,
			BranchTrue:  branchTrue,
			BranchFalse: branchFalse,
		}
		return
	}
	return
}

func parseLogicOr(tokens []Token, total, idx int) (index int, node Expression, err error) {
	index = idx
	var exp []Expression
	for token, ok := get(tokens, total, index); ok && (
	// logicAnd ( || logicAnd )*	index==idx 时是第一个 logicAnd
	index == idx || token.Value == "||"); token, ok = get(tokens, total, index) {
		if token.Value == "||" {
			_, index, err = consume(tokens, total, index, TokenTypeLGC, "||")
			if err != nil {
				return
			}
		}
		index, node, err = parseLogicAnd(tokens, total, index)
		if err != nil {
			return
		}
		exp = append(exp, node)
	}
	node = &LogicNode{Op: "||", Exps: exp}
	return
}

func parseLogicAnd(tokens []Token, total, idx int) (index int, node Expression, err error) {
	index = idx
	var exp []Expression
	for token, ok := get(tokens, total, index); ok && (
	// equality ( && equality )*	index==idx 时是第一个 equality
	index == idx || token.Value == "&&"); token, ok = get(tokens, total, index) {
		if token.Value == "&&" {
			_, index, err = consume(tokens, total, index, TokenTypeLGC, "&&")
			if err != nil {
				return
			}
		}
		index, node, err = parseEquality(tokens, total, index)
		if err != nil {
			return
		}
		exp = append(exp, node)
	}
	node = &LogicNode{Op: "&&", Exps: exp}
	return
}

func parseEquality(tokens []Token, total, idx int) (index int, node Expression, err error) {
	index = idx
	index, node, err = parseRelational(tokens, total, index)
	if err != nil {
		return
	}
	if token, ok := get(tokens, total, index); ok && token.Type == TokenTypeEQU {
		_, index, err = consume(tokens, total, index, TokenTypeEQU, "")
		if err != nil {
			return
		}
		ret := &CompareNode{
			Exp: node,
			Op:  token.Value,
		}
		node = ret
		index, ret.Other, err = parseRelational(tokens, total, index)
		if err != nil {
			return
		}
	}
	return
}

func parseRelational(tokens []Token, total, idx int) (index int, node Expression, err error) {
	index = idx
	index, node, err = parseAdd(tokens, total, index)
	if err != nil {
		return
	}
	if token, ok := get(tokens, total, index); ok && token.Type == TokenTypeCMP {
		_, index, err = consume(tokens, total, index, TokenTypeCMP, "")
		if err != nil {
			return
		}
		ret := &CompareNode{
			Exp: node,
			Op:  token.Value,
		}
		node = ret
		index, ret.Other, err = parseAdd(tokens, total, index)
		if err != nil {
			return
		}
	}
	return
}

func parseAdd(tokens []Token, total, idx int) (index int, node Expression, err error) {
	index = idx
	var exp Expression
	index, exp, err = parseMul(tokens, total, index)
	if err != nil {
		return
	}
	var op []string
	var other []Expression
	for token, ok := get(tokens, total, index); ok && token.Type == TokenTypeADD; token, ok = get(tokens, total, index) {
		_, index, err = consume(tokens, total, index, TokenTypeADD, "")
		if err != nil {
			return
		}
		index, node, err = parseMul(tokens, total, index)
		if err != nil {
			return
		}
		op = append(op, token.Value)
		other = append(other, node)
	}
	node = &BinaryNExp{
		Exp:   exp,
		Op:    op,
		Other: other,
	}
	return
}

func parseMul(tokens []Token, total, idx int) (index int, node Expression, err error) {
	index = idx
	var exp Expression
	index, exp, err = parseUnary(tokens, total, index)
	if err != nil {
		return
	}
	var op []string
	var other []Expression
	for token, ok := get(tokens, total, index); ok && token.Type == TokenTypeMUL; token, ok = get(tokens, total, index) {
		_, index, err = consume(tokens, total, index, TokenTypeMUL, "")
		if err != nil {
			return
		}
		index, node, err = parseUnary(tokens, total, index)
		if err != nil {
			return
		}
		op = append(op, token.Value)
		other = append(other, node)
	}
	node = &BinaryNExp{
		Exp:   exp,
		Op:    op,
		Other: other,
	}
	return
}

func parseUnary(tokens []Token, total, idx int) (index int, node Expression, err error) {
	index = idx
	op := ""
	if token, ok := get(tokens, total, index); ok && token.Value == "!" {
		_, index, err = consume(tokens, total, index, TokenTypeLGC, "!")
		if err != nil {
			return
		}
		op = "!"
	}
	index, node, err = parsePrimary(tokens, total, index)
	if err != nil {
		return
	}
	node = &UnaryExp{Op: op, Exp: node}
	return
}

func parsePrimary(tokens []Token, total, idx int) (index int, node Expression, err error) {
	index = idx
	if token, ok := get(tokens, total, index); ok {
		switch token.Type {
		case TokenTypeIDN:
			_, index, err = consume(tokens, total, index, TokenTypeIDN, "n")
			if err != nil {
				return
			}
			node = &PrimaryNode{Type: token.Type}
			return
		case TokenTypeNUM:
			_, index, err = consume(tokens, total, index, TokenTypeNUM, "")
			if err != nil {
				return
			}
			node = &PrimaryNode{Type: token.Type, Num: token.Number}
			return
		case TokenTypeLPA:
			_, index, err = consume(tokens, total, index, TokenTypeLPA, "(")
			if err != nil {
				return
			}
			index, node, err = parseExpression(tokens, total, index)
			if err != nil {
				return
			}
			_, index, err = consume(tokens, total, index, TokenTypeRPA, ")")
			if err != nil {
				return
			}
			node = &PrimaryNode{Type: token.Type, Exp: node}
			return
		}
		err = fmt.Errorf("expected ID, NUM, or '(', but got token %v", token)
		return
	}
	err = fmt.Errorf("expected ID, NUM, or '(' after token %v", tokens[total-1])
	return
}
