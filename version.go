package main

import (
	"fmt"
	"os"
)

const version = "0.1.0"
const nameAndVersion = "fserve v"+version

func printVersion() {
	fmt.Println(nameAndVersion)
	os.Exit(0)
}
