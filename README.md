# gonew
Simple tool to generate a new go project.

## Install
```bash
$ go get github.com/jonathantorres/gonew
```

## Creating a new project
This command will create your project inside in the location that you are currently in, by default it will create the `myproject` folder and also add a `main.go` file inside of it for your program.
```bash
$ gonew myproject
```

You can also use the `--gitignore` option to also generate a default `.gitignore` file
```bash
$ gonew --gitignore myproject
```
