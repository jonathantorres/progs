# ping
Test the reachability of a host via IP

## Install
Install by running the following command
```bash
$ go get github.com/jonathantorres/net/ping
```

## Usage and options
`-c` Stop after sending -c packets
```bash
$ ping -c 10 myurl.com
```

`-i` Wait -i seconds between sending each packet
```bash
$ ping -i 2 myurl.com
```

`-o` Exit successfully after receiving one reply packet
```bash
$ ping -o myurl.com
```

`-s` Specify the number of data bytes to be sent
```bash
$ ping -s 64 myurl.com
```

`-t` Timeout, in seconds before `ping` exits regardless of how many packets have been received
```bash
$ ping -t 10 myurl.com
```

`-4` Use IPv4 only
```bash
$ ping -4 myurl.com
```

`-6` Use IPv6 only
```bash
$ ping -6 myurl.com
```

`-d` Set the SO_DEBUG option on the socket being used
```bash
$ ping -d myurl.com
```
