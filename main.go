package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"unicode"
)

const maxVal = 100

func main() {
	r := bufio.NewReader(os.Stdin)
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("%s\n", err)
			continue
		}
		res, err := calc(line)
		if err != nil {
			fmt.Printf("%s\n", err)
			continue
		}
		fmt.Printf("\t%.8g\n", res)
	}
}

type stack struct {
	val []float64
	p   int
}

func newStack() *stack {
	v := make([]float64, maxVal)
	return &stack{
		val: v,
	}
}

func (s *stack) pop() float64 {
	if s.p > 0 {
		s.p--
		v := s.val[s.p]
		return v
	}
	fmt.Printf("error: stack empty\n")
	return 0.0
}

func (s *stack) push(f float64) error {
	if s.p < maxVal {
		s.val[s.p] = f
		s.p++
		return nil
	}
	return fmt.Errorf("error: stack full, can't push %f", f)
}

func calc(line []byte) (float64, error) {
	var num []byte
	st := newStack()
	for i, b := range line {
		c := rune(b)
		var isNegative bool
		if unicode.IsSpace(c) {
			if c != '\n' && len(num) > 0 {
				n, err := strconv.ParseFloat(string(num), 64)
				if err != nil {
					return 0.0, err
				}
				err = st.push(n)
				if err != nil {
					return 0.0, err
				}
				num = nil
			}
			continue
		}
		if c == '-' && i+1 < len(line) && unicode.IsDigit(rune(line[i+1])) {
			isNegative = true
		}
		if unicode.IsDigit(c) || c == '.' || isNegative {
			num = append(num, b)
			continue
		}
		switch c {
		case '+':
			st.push(st.pop() + st.pop())
			break
		case '*':
			st.push(st.pop() * st.pop())
			break
		case '-':
			op2 := st.pop()
			st.push(st.pop() - op2)
			break
		case '/':
			op2 := st.pop()
			if op2 != 0.0 {
				st.push(st.pop() / op2)
			} else {
				return 0.0, errors.New("error: zero divisor")
			}
			break
		case '%':
			op2 := st.pop()
			if op2 != 0.0 {
				st.push(math.Mod(st.pop(), op2))
			}
			break
		default:
			return 0.0, fmt.Errorf("error: unknown command %c", c)
		}
	}
	return st.pop(), nil
}
