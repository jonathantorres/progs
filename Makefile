ping: ping/main.go
	go build -o ping github.com/jonathantorres/net/ping

traceroute: traceroute/main.go
	go build -o traceroute github.com/jonathantorres/net/traceroute

static: static/main.go
	go build -o static github.com/jonathantorres/net/static

ftp: ftp/main.go
	go build -o ftp github.com/jonathantorres/net/ftp

# Run tests
.PHONY: test
test:
	go test .

.PHONY: clean
clean:
	go clean
