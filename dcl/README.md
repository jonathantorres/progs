# dcl
Parse a C declaration into words

## Install
Install by running the following command
```bash
$ go get github.com/jonathantorres/cmd/dcl
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
```
