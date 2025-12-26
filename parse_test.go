package plurals

import (
	"fmt"
	"testing"
)

func TestCompile(t *testing.T) {
	for _, tt := range []struct {
		exp string
		err bool
	}{
		{exp: "", err: true},
		{exp: "0", err: false},
		{exp: "n != 1", err: false},
		{exp: "n % 10 == 1 && n % 100 != 11 ? 0 : n != 0 ? 1 : 2", err: false},
		{exp: "n == 1 ? 0 : n == 2 ? 1 : 2", err: false},
		{exp: "n == 1 ? 0 : ( n == 0 || ( n % 100 > 0 && n % 100 < 20 ) ) ? 1 : 2", err: false},
		{exp: "n % 10 == 1 && n % 100 != 11 ? 0 : n % 10 >= 2 && ( n % 100 < 10 || n % 100 >= 20 ) ? 1 : 2", err: false},
		{exp: "n % 10 == 1 && n % 100 != 11 ? 0 : n % 10 >= 2 && n % 10 <= 4 && ( n % 100 < 10 || n % 100 >= 20 ) ? 1 : 2", err: false},
		{exp: "( n == 1 ) ? 0 : ( n >= 2 && n <= 4 ) ? 1 : 2", err: false},
		{exp: "n == 1 ? 0 : n % 10 >= 2 && n % 10 <= 4 && ( n % 100 < 10 || n % 100 >= 20 ) ? 1 : 2", err: false},
		{exp: "n % 100 == 1 ? 0 : n % 100 == 2 ? 1 : n % 100 == 3 || n % 100 == 4 ? 2 : 3", err: false},
		{exp: "n == 0 ? 0 : n == 1 ? 1 : n == 2 ? 2 : n % 100 >= 3 && n % 100 <= 10 ? 3 : n % 100 >= 11 ? 4 : 5", err: false},
		{exp: "n==;n", err: true},
		{exp: "1>!", err: true},
		{exp: "! == 1", err: true},
	} {
		exp, err := Compile(tt.exp)
		t.Logf("%q: exp=%v, err=%+v", tt.exp, exp, err)
		if tt.err != (err != nil) {
			t.Errorf("fail want err=%v", tt.err)
			continue
		}
		if err == nil && fmt.Sprintf("%v", exp) != tt.exp {
			t.Errorf("got: %v, want: %s", exp, tt.exp)
			continue
		}
		if err == nil {
			for n := range int64(3) {
				got, err := Eval(tt.exp, n)
				t.Logf("Eval(%v)=%v, err=%v", n, got, err)
			}
		}
	}
}

func TestEval(t *testing.T) {
	for s, f := range commons {
		for n := range 1000 {
			n := int64(n)
			val, err := eval(s, n)
			if err != nil {
				t.Errorf("exp=%q, err=%+v", s, err)
				break
			}
			want := f(n)
			if val != want {
				t.Errorf("Eval fail: n=%d, got=%v, want=%v, %q",
					n, val, want, s)
				break
			}
		}
	}
}
