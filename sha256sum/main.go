package main

import (
	"bufio"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"os"
)

var use384 = flag.Bool("a", false, "Use SHA384")
var use512 = flag.Bool("b", false, "Use SHA512")

func main() {
	flag.Parse()
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		b := s.Bytes()
		if *use384 {
			fmt.Printf("%x\n", sha512.Sum384(b))
		} else if *use512 {
			fmt.Printf("%x\n", sha512.Sum512(b))
		} else {
			fmt.Printf("%x\n", sha256.Sum256(b))
		}
	}
	if err := s.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error! %s\n", err)
	}
}
