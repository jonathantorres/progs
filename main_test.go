package main

import "testing"

func TestDeclarations(t *testing.T) {
	tests := []struct {
		exp  []byte
		want string
	}{
		{[]byte("char **argv"), "argv: pointer to pointer to char"},
		{[]byte("int (*daytab)[13]"), "daytab: pointer to array[13] of int"},
		{[]byte("int *daytab[13]"), "daytab: array[13] of pointer to int"},
		{[]byte("void *comp()"), "comp: function returning pointer to void"},
		{[]byte("void (*comp)()"), "comp: pointer to function returning void"},
		{[]byte("char (*(*x())[])()"), "x: function returning pointer to array[] of pointer to function returning char"},
		{[]byte("char (*(*x[3])())[5]"), "x: array[3] of pointer to function returning pointer to array[5] of char"},
	}
	for _, test := range tests {
		got, err := parse(test.exp)
		if err != nil {
			t.Errorf("%s", err)
		}
		if got != test.want {
			t.Errorf("expected: %s, got: %s", test.want, got)
		}
	}
}
