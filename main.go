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
	var printVersion = flag.Bool("v", false, "Print gonew's version")
	var name string
	flag.Parse()
	if *printVersion {
		fmt.Fprintf(os.Stdout, "gonew v%s\n", version)
		os.Exit(0)
	}
	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "usage: gonew [-g -r -p -v] [project name]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	name = flag.Args()[0]
	createRootFolder(&name)
	if *isPackage {
		if err := createPackageFile(&name); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		if err := createMain(&name); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}
	if *gitignore {
		if err := createGitIgnore(&name); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}
	if *readme {
		if err := createReadme(&name); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func createPackageFile(path *string) error {
	packageStr := strings.Replace(packageText, "{project}", *path, -1)
	err := createFile(*path+"/"+*path+".go", []byte(packageStr))
	if err != nil {
		return fmt.Errorf("there was a problem creating %s.go file, %s\n", *path, err)
	}
	return nil
}

func createReadme(path *string) error {
	readme := strings.Replace(readmeText, "{project}", *path, -1)
	err := createFile(*path+"/README.md", []byte(readme))
	if err != nil {
		fmt.Errorf("there was a problem creating README.md file, %s\n", err)
	}
	return nil
}

func createRootFolder(name *string) error {
	if err := os.Mkdir(*name, 0777); err != nil {
		return fmt.Errorf("directory %s already exists\n", *name)
	}
	return nil
}

func createGitIgnore(path *string) error {
	err := createFile(*path+"/.gitignore", []byte(gitignoreText))
	if err != nil {
		return fmt.Errorf("there was a problem creating .gitignore file, %s\n", err)
	}
	return nil
}

func createMain(path *string) error {
	err := createFile(*path+"/main.go", []byte(mainText))
	if err != nil {
		return fmt.Errorf("there was a problem creating main.go file, %s\n", err)
	}
	return nil
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
