PROG := ftp

# compile program
$(PROG): main.go
	go build -o ftp github.com/jonathantorres/net/ftp

# Run tests
.PHONY: test
test:
	go test .

.PHONY: clean
clean:
	go clean
