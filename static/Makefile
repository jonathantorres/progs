PROG := static

# compile program
$(PROG): main.go
	go build -o static github.com/jonathantorres/net/static

# Run tests
.PHONY: test
test:
	go test .

.PHONY: clean
clean:
	go clean

release:
	rm -fr ./bin && mkdir ./bin && touch ./bin/.gitkeep
	GOOS=darwin GOARCH=amd64 go build && mv ./static ./bin/static_darwin
	GOOS=freebsd GOARCH=amd64 go build && mv ./static ./bin/static_freebsd
	GOOS=linux GOARCH=amd64 go build && mv ./static ./bin/static_linux
	GOOS=netbsd GOARCH=amd64 go build && mv ./static ./bin/static_netbsd
	GOOS=plan9 GOARCH=amd64 go build && mv ./static ./bin/static_plan9
	GOOS=windows GOARCH=amd64 go build && mv ./static.exe ./bin/static_windows.exe
