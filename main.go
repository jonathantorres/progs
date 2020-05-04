package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("usage: gonew [project]")
		os.Exit(1)
	}
	var files []string = []string{
		gitignoreText,
		mainText,
		readmeText,
	}
	for _, file := range files {
		fmt.Println(file)
	}
}
