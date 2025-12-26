package plurals

import (
	"reflect"
	"testing"
)

func TestLex(t *testing.T) {
	for _, tt := range []struct {
		s      string
		tokens []string
		err    bool
	}{
		{s: "", tokens: []string{}, err: false},
		{s: " ", tokens: []string{}, err: false},
		{s: "0", tokens: []string{"0"}, err: false},
		{s: " 0", tokens: []string{"0"}, err: false},
		{s: " 10 ", tokens: []string{"10"}, err: false},
		{s: "n", tokens: []string{"n"}, err: false},
		{s: "==", tokens: []string{"=="}, err: false},
		{s: "=", tokens: []string{}, err: true},
		{s: "!=", tokens: []string{"!="}, err: false},
		{s: "!", tokens: []string{"!"}, err: false},
		{s: "&&", tokens: []string{"&&"}, err: false},
		{s: "||", tokens: []string{"||"}, err: false},
		{s: "|", tokens: []string{}, err: true},
		{s: "&", tokens: []string{}, err: true},
		{s: ">", tokens: []string{">"}, err: false},
		{s: ">=", tokens: []string{">="}, err: false},
		{s: "<", tokens: []string{"<"}, err: false},
		{s: "<=", tokens: []string{"<="}, err: false},
		{s: "+", tokens: []string{"+"}, err: false},
		{s: "-", tokens: []string{"-"}, err: false},
		{s: "*", tokens: []string{"*"}, err: false},
		{s: "/", tokens: []string{"/"}, err: false},
		{s: "%", tokens: []string{"%"}, err: false},
		{s: "?", tokens: []string{"?"}, err: false},
		{s: ":", tokens: []string{":"}, err: false},
		{s: "(", tokens: []string{"("}, err: false},
		{s: ")", tokens: []string{")"}, err: false},
		{s: "a", tokens: []string{}, err: true},
	} {
		tokens, err := Lex(tt.s)
		t.Logf("tokens=%v, err=%v", tokens, err)
		var got = make([]string, 0, len(tokens))
		for _, tk := range tokens {
			got = append(got, tk.Value)
		}
		if tt.err != (err != nil) {
			t.Errorf("Fail: `%s`, err=%+v", tt.s, err)
		}
		if !reflect.DeepEqual(got, tt.tokens) {
			t.Errorf("got=%v, want=%v", got, tt.tokens)
		}
	}

}
