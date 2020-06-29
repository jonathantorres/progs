package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var gitignore = flag.Bool("g", false, "Generate a .gitignore file")
	var readme = flag.Bool("r", false, "Generate a README.md file")
	var isPackage = flag.Bool("p", false, "Generate a package file, instead of a binary")
	var name string
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "usage: gonew [-g -r -p] [project name]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	name = flag.Args()[0]
	createRootFolder(&name)
	if *isPackage {
		createPackageFile(&name)
	} else {
		createMain(&name)
	}
	if *gitignore {
		createGitIgnore(&name)
	}
	if *readme {
		createReadme(&name)
	}
}

func createPackageFile(path *string) {

}

func createReadme(path *string) {

}

func createRootFolder(name *string) {
	if err := os.Mkdir(*name, 0777); err != nil {
		fmt.Println("directory " + *name + " already exists!")
	}
}

func createGitIgnore(path *string) {
	file, err := os.Create(*path + "/.gitignore")
	if err != nil {
		fmt.Println("Could not create .gitignore file:", err)
		return
	}
	defer file.Close()
	_, err = file.Write([]byte(gitignoreText))
	if err != nil {
		fmt.Println("Could not write code to .gitignore file:", err)
		return
	}
}

func createMain(path *string) {
	file, err := os.Create(*path + "/main.go")
	if err != nil {
		fmt.Println("Could not create main.go file:", err)
		return
	}
	defer file.Close()
	_, err = file.Write([]byte(mainText))
	if err != nil {
		fmt.Println("Could not write code to main.go file:", err)
		return
	}
}
