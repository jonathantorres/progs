package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

const (
	name rune = iota
	parens
	brackets
	syntaxErr
)

var buf *bytes.Buffer
var token bytes.Buffer
var out bytes.Buffer
var iname string
var datatype string
var tokentype rune

func main() {
	r := bufio.NewReader(os.Stdin)
	for {
		l, err := r.ReadBytes('\n')
		if err != nil {
			break
		}
		decl, err := parse(l)
		if err != nil {
			fmt.Printf("%s\n", err)
			continue
		}
		fmt.Printf("%s\n", decl)
	}
}

func parse(exp []byte) (string, error) {
	buf = bytes.NewBuffer(exp)
	for {
		_, err := getToken()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		datatype = token.String()
		token.Reset()
		err = dcl()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		if tokentype == ']' {
			return "", errors.New("syntax error, missing ]")
		}
	}
	outStr := out.String()
	outStr = strings.TrimSpace(cleanOutStr(outStr))
	s := fmt.Sprintf("%s: %s %s", iname, outStr, datatype)
	out.Reset()

	return s, nil
}

func getToken() (rune, error) {
	var c rune
	var err error
	for {
		c, _, err = buf.ReadRune()
		if err != nil {
			return 0, err
		}
		if unicode.IsSpace(c) {
			continue
		}
		if c == '(' {
			c, _, err := buf.ReadRune()
			if err != nil {
				if err == io.EOF {
					return 0, nil
				}
				return 0, err
			}
			if c == ')' {
				tokentype = parens
				return tokentype, nil
			}
			buf.UnreadRune()
			tokentype = '('
			return tokentype, nil
		} else if c == '[' {
			last := writeTokenChars(c)
			if last == ']' {
				tokentype = brackets
				return tokentype, nil
			}
			tokentype = ']'
			return tokentype, nil
		} else if unicode.IsLetter(c) {
			writeTokenName(c)
			tokentype = name
			return tokentype, nil
		} else {
			tokentype = c
			return tokentype, nil
		}
	}
	return 0, errors.New("token not found")
}

func dcl() error {
	var ns int
	for {
		r, err := getToken()
		if err != nil {
			return err
		}
		if r == '*' {
			ns++
		} else {
			break
		}
	}
	err := dirdcl()
	if err != nil && err != io.EOF {
		return err
	}
	for i := ns; i > 0; i-- {
		out.WriteString("pointer to ")
	}
	return nil
}

func dirdcl() error {
	if tokentype == '(' {
		err := dcl()
		if err != nil {
			return err
		}
		if tokentype != ')' {
			err := errors.New("error: missing ) syntax error")
			tokentype = syntaxErr
			return err
		}
	} else if tokentype == name {
		iname = token.String()
		token.Reset()
	} else {
		err := errors.New("error: expected name or (dcl) syntax error")
		tokentype = syntaxErr
		return err
	}
	for {
		r, err := getToken()
		if err != nil {
			return err
		}
		if r == parens || r == brackets {
			if r == parens {
				out.WriteString("function returning ")
			} else {
				out.WriteString("array")
				out.WriteString(token.String())
				out.WriteString(" of ")
				token.Reset()
			}
		} else {
			break
		}
	}
	return nil
}

func writeTokenChars(f rune) rune {
	var last rune
	token.WriteRune(f) // write first rune
	for {
		c, _, err := buf.ReadRune()
		if err != nil {
			break
		}
		// continue writing until we find the closing ]
		if c == ']' || c == '\n' {
			if c == ']' {
				token.WriteRune(c)
				last = c
			}
			break
		}
		token.WriteRune(c)
		last = c
	}
	return last
}

func writeTokenName(f rune) {
	token.WriteRune(f) // write first rune in name
	for {
		c, _, err := buf.ReadRune()
		if err != nil {
			break
		}
		// write the rest of the name
		if unicode.IsLetter(c) || unicode.IsNumber(c) {
			token.WriteRune(c)
			continue
		}
		// stop here, since we found another character
		// push it back into the buffer
		buf.UnreadRune()
		break
	}
}

func cleanOutStr(outStr string) string {
	var res strings.Builder
	ss := strings.Split(outStr, " ")
	for i, s := range ss {
		res.WriteString(strings.TrimSpace(s))
		if i < len(ss)-1 {
			res.WriteString(" ")
		}
	}
	return res.String()
}
