PROG := traceroute

# compile program
$(PROG): main.go
	go build -o traceroute github.com/jonathantorres/net/traceroute

# Run tests
.PHONY: test
test:
	go test .

.PHONY: clean
clean:
	go clean
