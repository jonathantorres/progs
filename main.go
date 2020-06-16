package main

import (
	"fmt"
	"os"
	"flag"
)

func main() {
	var gitignore = flag.Bool("gitignore", false, "Generate a .gitignore file")
	var name string
	flag.Parse()
	if (len(flag.Args()) == 0) {
		fmt.Println("usage: gonew [--options] [project_name]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	name = flag.Args()[0]
	createRootFolder(&name)
	createMain(&name)
	if (*gitignore) {
		createGitIgnore(&name);
	}
}

func createRootFolder(name *string) {
	if err := os.Mkdir(*name, 0777); err != nil {
		fmt.Println("directory " + *name + " already exists!")
	}
}

func createGitIgnore(path *string) {
	file, err := os.Create(*path + "/.gitignore")
	if (err != nil) {
		fmt.Println("Could not create .gitignore file:", err)
		return
	}
	_, err = file.Write([]byte(gitignoreText))
	if (err != nil) {
		fmt.Println("Could not write code to .gitignore file:", err)
		return
	}
}

func createMain(path *string) {
	file, err := os.Create(*path + "/main.go")
	if (err != nil) {
		fmt.Println("Could not create main.go file:", err)
		return
	}
	_, err = file.Write([]byte(mainText))
	if (err != nil) {
		fmt.Println("Could not write code to main.go file:", err)
		return
	}
}
