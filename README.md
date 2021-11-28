# zing
Just like the good old PING, but with a Z.

## Install
Install by running the following command
```bash
$ go get github.com/jonathantorres/zing
```

## Usage and options
`-c` Stop after sending -c packets
```bash
$ zing -c 10 myurl.com
```

`-i` Wait -i seconds between sending each packet
```bash
$ zing -i 2 myurl.com
```

`-o` Exit successfully after receiving one reply packet
```bash
$ zing -o myurl.com
```

`-s` Specify the number of data bytes to be sent
```bash
$ zing -s 64 myurl.com
```

`-t` Timeout, in seconds before `zing` exits regardless of how many packets have been received
```bash
$ zing -t 10 myurl.com
```

`-d` Set the SO_DEBUG option on the socket being used
```bash
$ zing -d myurl.com
```
