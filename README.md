# gonew
Simple tool to generate a new go project.

## Install
If you have go installed, you can install by running the following command, otherwise you can download the binary for your operating system from the latest release [here](https://github.com/jonathantorres/gonew/releases/tag/v0.1.0).
```bash
$ go get github.com/jonathantorres/gonew
```

## Creating a new project
This command will create your project inside in the location that you are currently in, by default it will create the `myproject` folder and also add a `main.go` file inside of it for your program.
```bash
$ gonew myproject
```

The following options are also available:

`-g` will generate a default `.gitignore` file
```bash
$ gonew  -g myproject
```

`-r` will generate a `README` file with your project's name
```bash
$ gonew  -r myproject
```

`-p` will generate a file `myproject.go` instead of a `main.go` file. Use this option if you are generating a package instead of an executable
```bash
$ gonew  -p myproject
```
