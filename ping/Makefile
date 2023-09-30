PROG := ping

# compile program
$(PROG): main.go
	go build -o ping github.com/jonathantorres/net/ping

# Run tests
.PHONY: test
test:
	go test .

.PHONY: clean
clean:
	go clean
