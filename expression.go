package plurals

import (
	"fmt"
	"strings"
)

type Expression interface {
	Eval(n int64) (int64, error)
}

func Eval(s string, n int64) (int64, error) {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\t", "")
	if f, ok := commons[s]; ok {
		return f(n), nil
	}
	return eval(s, n)
}

func eval(s string, n int64) (int64, error) {
	if exp, ok := cache[s]; ok {
		return exp.Eval(n)
	}
	exp, err := Compile(s)
	if err != nil {
		return 0, err
	}
	cache[s] = exp
	return exp.Eval(n)
}

var cache = map[string]Expression{}

const (
	nFalse = 0
	nTrue  = 1
)

func b2i(b bool) int64 {
	if b {
		return nTrue
	}
	return nFalse
}

func i2b(i int64) bool {
	return i != nFalse
}

type TernaryNode struct {
	Condition   Expression
	BranchTrue  Expression
	BranchFalse Expression
}

func (e *TernaryNode) Eval(n int64) (int64, error) {
	val, err := e.Condition.Eval(n)
	if err != nil {
		return 0, err
	}
	if i2b(val) {
		return e.BranchTrue.Eval(n)
	}
	return e.BranchFalse.Eval(n)
}

func (e *TernaryNode) String() string {
	return fmt.Sprintf("%v ? %v : %v", e.Condition, e.BranchTrue, e.BranchFalse)
}

type LogicNode struct {
	Op   string
	Exps []Expression
}

func (e *LogicNode) Eval(n int64) (val int64, err error) {
	op := e.Op
	for i, exp := range e.Exps {
		if i == 0 {
			val, err = exp.Eval(n)
			if err != nil {
				return 0, err
			}
			continue
		}
		switch op {
		case "&&":
			i, err := exp.Eval(n)
			if err != nil {
				return 0, err
			}
			val = b2i(i2b(val) && i2b(i))
		case "||":
			i, err := exp.Eval(n)
			if err != nil {
				return 0, err
			}
			val = b2i(i2b(val) || i2b(i))
		default:
			return 0, fmt.Errorf("assert failed")
		}
		if !i2b(val) && op == "&&" {
			return nFalse, nil
		}
		if i2b(val) && op == "||" {
			return nTrue, nil
		}
	}
	return val, nil
}

func (e *LogicNode) String() string {
	var sb strings.Builder
	for i, exp := range e.Exps {
		if i > 0 {
			fmt.Fprintf(&sb, " %v ", e.Op)
		}
		fmt.Fprintf(&sb, "%v", exp)
	}
	return sb.String()
}

type CompareNode struct {
	Exp   Expression
	Op    string
	Other Expression
}

func (e *CompareNode) Eval(n int64) (int64, error) {
	val, err := e.Exp.Eval(n)
	if err != nil {
		return 0, err
	}
	if e.Other != nil {
		i, err := e.Other.Eval(n)
		if err != nil {
			return 0, err
		}
		switch e.Op {
		case "==":
			return b2i(val == i), nil
		case "!=":
			return b2i(val != i), nil
		case ">":
			return b2i(val > i), nil
		case ">=":
			return b2i(val >= i), nil
		case "<":
			return b2i(val < i), nil
		case "<=":
			return b2i(val <= i), nil
		default:
			return 0, fmt.Errorf("assert failed")
		}
	}
	return val, nil
}

func (e *CompareNode) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%v", e.Exp)
	if e.Other != nil {
		fmt.Fprintf(&sb, " %v %v", e.Op, e.Other)
	}
	return sb.String()
}

type BinaryNExp struct {
	Exp   Expression
	Op    []string
	Other []Expression
}

func (e *BinaryNExp) Eval(n int64) (int64, error) {
	val, err := e.Exp.Eval(n)
	if err != nil {
		return 0, err
	}
	for idx, other := range e.Other {
		i, err := other.Eval(n)
		if err != nil {
			return 0, err
		}
		switch e.Op[idx] {
		case "+":
			val = val + i
		case "-":
			val = val - i
		case "*":
			val = val * i
		case "/":
			if i == 0 {
				return 0, fmt.Errorf("divide zero")
			}
			val = val / i
		case "%":
			if i == 0 {
				return 0, fmt.Errorf("divide zero")
			}
			val = val % i
		default:
			return 0, fmt.Errorf("assert failed")
		}
	}
	return val, nil
}

func (e *BinaryNExp) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%v", e.Exp)
	for i, other := range e.Other {
		fmt.Fprintf(&sb, " %v %v", e.Op[i], other)
	}
	return sb.String()
}

type UnaryExp struct {
	Op  string
	Exp Expression
}

func (e *UnaryExp) Eval(n int64) (int64, error) {
	val, err := e.Exp.Eval(n)
	if err != nil {
		return 0, err
	}
	if e.Op == "!" {
		// 取反 !(n==1)
		return b2i(!i2b(val)), nil
	}
	return val, nil
}

func (e *UnaryExp) String() string {
	if e.Op == "!" {
		return fmt.Sprintf("!%v", e.Exp)
	}
	return fmt.Sprintf("%v", e.Exp)
}

type PrimaryNode struct {
	Type TokenType
	Num  int64
	Exp  Expression
}

func (e *PrimaryNode) Eval(n int64) (int64, error) {
	switch e.Type {
	case TokenTypeIDN:
		return n, nil
	case TokenTypeNUM:
		return e.Num, nil
	case TokenTypeLPA:
		return e.Exp.Eval(n)
	}
	return 0, fmt.Errorf("assert failed")
}

func (e *PrimaryNode) String() string {
	switch e.Type {
	case TokenTypeIDN:
		return "n"
	case TokenTypeNUM:
		return fmt.Sprintf("%v", e.Num)
	case TokenTypeLPA:
		return fmt.Sprintf("( %v )", e.Exp)
	}
	return ""
}
