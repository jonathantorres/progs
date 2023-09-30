# traceroute
Display the possible route and measure transit delays of packets. It uses UDP as the probing mechanism.

## Install
Install by running the following command
```bash
$ go get github.com/jonathantorres/net/traceroute
```

## Usage and options
Default usage
```bash
$ traceroute myurl.com
```

`-d` Enable socket level debugging (if supported)
```bash
$ traceroute -d myurl.com
```

`-f` Specify with what TTL to start. Defaults to 1
```bash
$ traceroute -f 1 myurl.com
```

`-m` Specify the maximum number of hops (max time-to-live value) the program will probe. The default is 30
```bash
$ traceroute -m 30 myurl.com
```

`-p` Specify the destination port to use. This number will be incremented by each probe
```bash
$ traceroute -p 34500 myurl.com
```

`-q` Sets the number of probe packets per hop. The default number is 3
```bash
$ traceroute -q 3 myurl.com
```

`-4` Use IPv4 only
```bash
$ traceroute -4 myurl.com
```

`-6` Use IPv6 only
```bash
$ traceroute -6 myurl.com
```

`-w` Probe timeout. Specify how many seconds to wait for a response to a probe. Default value is 5
```bash
$ traceroute -w 5 myurl.com
```

`-z` Minimum amount of time to wait between probes (in seconds). The default is 0
```bash
$ traceroute -z 1 myurl.com
```
