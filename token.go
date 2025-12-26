package plurals

import "fmt"

type TokenType string

const (
	TokenTypeEOF = "EOF"
	TokenTypeERR = "ERR"
	TokenTypeNUM = "NUM" // \d+
	TokenTypeEQU = "EQU" // == !=
	TokenTypeLGC = "LGC" // || && !
	TokenTypeCMP = "CMP" // > >= < <=
	TokenTypeMUL = "MUL" // * / %
	TokenTypeADD = "ADD" // + -
	TokenTypeIDN = "IDN" // n
	TokenTypeQST = "QST" // ?
	TokenTypeCOL = "COL" // :
	TokenTypeLPA = "LPA" // (
	TokenTypeRPA = "RPA" // )
	TokenTypeCOM = "COM" // ;
)

type Token struct {
	Type   TokenType
	Value  string
	Number int64
	Start  int
	End    int
}

func (t Token) String() string {
	if t.Type == TokenTypeNUM {
		return fmt.Sprintf("{[%d:%d] %v(%v)}", t.Start, t.End, t.Type, t.Number)
	}
	return fmt.Sprintf("{[%d:%d] %v(%v)}", t.Start, t.End, t.Type, t.Value)
}

// https://www.gnu.org/software/gettext/manual/html_node/Plural-forms.html#index-specifying-plural-form-in-a-PO-file
// Plural-Forms: nplurals=2; plural=n == 1 ? 0 : 1;
// The nplurals value must be a decimal number which specifies how many different plural forms exist for this language.
// The string following plural is an expression which is using the C language syntax.
// Exceptions are that no negative numbers are allowed, numbers must be decimal, and the only variable allowed is n.
// Spaces are allowed in the expression, but backslash-newlines are not.

/*

exp : exp '?' exp ':' exp
    | exp '||' exp
    | exp '&&' exp
    | exp ('=='|'!=') exp
    | exp ('>'|'>='|'<'|'<=') exp
    | exp ('+'|'-') exp
    | exp ('*'|'/'|'%') exp
    | '!' exp
    | 'n'
    | NUMBER
    | '(' exp ')'
    ;

plural                    : expression ';'
expression                : ternary_expression
ternary_expression        : logical_or_expression ( '?' expression ':' expression )?
logical_or_expression     : logical_and_expression ( '||' logical_and_expression )*
logical_and_expression    : equality_expression ( '&&' equality_expression )*
equality_expression       : relational_expression ( ('=='|'!=') relational_expression )?
relational_expression     : additive_expression ( ('>'|'<'|'>='|'<=') additive_expression )?
additive_expression       : multiplicative_expression ( ('+'|'-') multiplicative_expression )*
multiplicative_expression : unary_expression ( ('*'|'/'|'%') unary_expression)*
unary_expression          : '!'? primary_expression
primary_expression        : 'n'
                          | NUMBER
                          | '(' expression ')'
*/
