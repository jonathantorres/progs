package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
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
	packageStr := strings.Replace(packageText, "{project}", *path, -1)
	err := createFile(*path+"/"+*path+".go", []byte(packageStr))
	if err != nil {
		fmt.Fprintf(os.Stderr, "there was a problem creating %s.go file, %s\n", *path, err)
	}
}

func createReadme(path *string) {
	readme := strings.Replace(readmeText, "{project}", *path, -1)
	err := createFile(*path+"/README.md", []byte(readme))
	if err != nil {
		fmt.Fprintf(os.Stderr, "there was a problem creating README.md file, %s\n", err)
	}
}

func createRootFolder(name *string) {
	if err := os.Mkdir(*name, 0777); err != nil {
		fmt.Fprintf(os.Stderr, "directory %s already exists\n", *name)
	}
}

func createGitIgnore(path *string) {
	err := createFile(*path+"/.gitignore", []byte(gitignoreText))
	if err != nil {
		fmt.Fprintf(os.Stderr, "there was a problem creating .gitignore file, %s\n", err)
	}
}

func createMain(path *string) {
	err := createFile(*path+"/main.go", []byte(mainText))
	if err != nil {
		fmt.Fprintf(os.Stderr, "there was a problem creating main.go file, %s\n", err)
	}
}

func createFile(path string, filedata []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(filedata)
	if err != nil {
		return err
	}
	return nil
}
