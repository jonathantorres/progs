PROG := zing

# compile program
$(PROG): main.go
	go build -o zing github.com/jonathantorres/zing

# Run tests
.PHONY: test
test:
	go test .

.PHONY: clean
clean:
	go clean
