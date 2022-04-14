PROG := rt

# compile program
$(PROG): main.go
	go build -o rt github.com/jonathantorres/rt

# Run tests
.PHONY: test
test:
	go test .

.PHONY: clean
clean:
	go clean
