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
=======
# dcl
Parse a C declaration into words
=======
# poli
A Reverse Polish calculator
>>>>>>> poli/master

## Install
Install by running the following command
```bash
<<<<<<< HEAD
$ go get github.com/jonathantorres/dcl
```

## Usage
After installing, run the program `dcl` and enter a C declaration. It will parse the declaration and give a word description of the declaration, here's an example session:
```bash
char **argv
argv: pointer to pointer to char

int (*daytab)[13]
daytab: pointer to array[13] of int

int *daytab[13]
daytab: array[13] of pointer to int

void *comp()
comp: function returning pointer to void

void (*comp)()
comp: pointer to function returning void
>>>>>>> dcl/master
=======
$ go get github.com/jonathantorres/poli
```

## Usage
After installing, run the program `poli` and start entering the commands to do the calculations, the basic operations are available (addition, substraction, multiplication, division, modulo)
```bash
100 100 +
        200

1 2 - 4 5 + *
        -9
```
