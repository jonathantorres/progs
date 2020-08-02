release:
	rm -fr ./bin && mkdir ./bin && touch ./bin/.gitkeep
	GOOS=darwin GOARCH=amd64 go build && mv ./fserve ./bin/fserve_darwin
	GOOS=freebsd GOARCH=amd64 go build && mv ./fserve ./bin/fserve_freebsd
	GOOS=linux GOARCH=amd64 go build && mv ./fserve ./bin/fserve_linux
	GOOS=netbsd GOARCH=amd64 go build && mv ./fserve ./bin/fserve_netbsd
	GOOS=plan9 GOARCH=amd64 go build && mv ./fserve ./bin/fserve_plan9
	GOOS=windows GOARCH=amd64 go build && mv ./fserve.exe ./bin/fserve_windows.exe
