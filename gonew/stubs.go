package main

var gitignoreText string = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with ` + "`go test -c`" + `
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/
`

var mainText string = `package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello world")
}
`

var packageText string = `package {project}

import (
	"fmt"
)

func init() {
	fmt.Println("hello world")
}
`

var readmeText string = `# {project}
`
