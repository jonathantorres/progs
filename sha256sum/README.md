# sha256sum
Produce the sha256 digest of data from standard input

## Install
Install by running the following command
```bash
$ go get github.com/jonathantorres/cmd/sha256sum
```

## Usage
After building the program, pass your data to it from stdin as usual. It will return the sha256 digest of whatever input is passed to it.
```bash
echo "this is some data" | ./sha256sum
```
