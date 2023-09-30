package main

import (
	"math/big"
	"testing"
)

func TestBasicCalculations(t *testing.T) {
	tests := []struct {
		op   []byte
		want float64
	}{
		{[]byte("100 100 +"), 200},
		{[]byte("1 2 - 4 5 + *"), -9},
		{[]byte("6434 800 -"), 5634},
		{[]byte("5 5 *"), 25},
		{[]byte("1000 100 /"), 10},
		{[]byte("700 56 %"), 28},
	}
	for _, test := range tests {
		got, err := calc(test.op)
		if err != nil {
			t.Fatalf("%s", err)
		}
		if big.NewFloat(got).Cmp(big.NewFloat(test.want)) != 0 {
			t.Errorf("calculation failed: got %f, but want %f", got, test.want)
		}
	}
}
