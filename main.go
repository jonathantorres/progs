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
	// var project string = os.Args[1]
	var files []string = []string{
		"main",
		"gitignore",
		"readme",
	}
	for _, file := range files {
		fmt.Println(file)
	}
}
