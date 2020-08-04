package main

import (
	"io/ioutil"
	"os"
	"testing"
	"strings"
)

var folderName = "test1"

func TestCreateRootFolder(t *testing.T) {
	if err := createRootFolder(&folderName); err != nil {
		t.Errorf("%s", err)
	}
	removeTestFolder()
}

func TestCreatePackageFile(t *testing.T) {
	if err := createRootFolder(&folderName); err != nil {
		t.Errorf("%s", err)
	}
	if err := createPackageFile(&folderName); err != nil {
		t.Errorf("%s", err)
	}
	fileContents, _ := ioutil.ReadFile(folderName+"/"+folderName+".go")
	stubContents := strings.Replace(packageText, "{project}", folderName, -1)
	if string(fileContents) != stubContents {
		t.Errorf("generated file %s is not the same as stub", folderName+"/"+folderName+".go")
		t.Logf("%s", string(fileContents))
		t.Logf("%s", stubContents)
	}
	removeTestFolder()
}

func TestCreateMainFile(t *testing.T) {
	if err := createRootFolder(&folderName); err != nil {
		t.Errorf("%s", err)
	}
	if err := createMain(&folderName); err != nil {
		t.Errorf("%s", err)
	}
	fileContents, _ := ioutil.ReadFile(folderName+"/"+"main.go")
	stubContents := mainText
	if string(fileContents) != stubContents {
		t.Errorf("generated file %s is not the same as stub", folderName+"/main.go")
		t.Logf("%s", string(fileContents))
		t.Logf("%s", stubContents)
	}
	removeTestFolder()
}

func TestCreateGitIgnore(t *testing.T) {
	if err := createRootFolder(&folderName); err != nil {
		t.Errorf("%s", err)
	}
	if err := createGitIgnore(&folderName); err != nil {
		t.Errorf("%s", err)
	}
	fileContents, _ := ioutil.ReadFile(folderName+"/"+".gitignore")
	stubContents := gitignoreText
	if string(fileContents) != stubContents {
		t.Errorf("generated file %s is not the same as stub", folderName+"/.gitignore")
		t.Logf("%s", string(fileContents))
		t.Logf("%s", stubContents)
	}
	removeTestFolder()
}

func TestCreateReadme(t *testing.T) {
	if err := createRootFolder(&folderName); err != nil {
		t.Errorf("%s", err)
	}
	if err := createReadme(&folderName); err != nil {
		t.Errorf("%s", err)
	}
	fileContents, _ := ioutil.ReadFile(folderName+"/README.md")
	stubContents := strings.Replace(readmeText, "{project}", folderName, -1)
	if string(fileContents) != stubContents {
		t.Errorf("generated file %s is not the same as stub", folderName+"/README.md")
		t.Logf("%s", string(fileContents))
		t.Logf("%s", stubContents)
	}
	removeTestFolder()
}

func removeTestFolder() {
	os.RemoveAll(folderName)
}
