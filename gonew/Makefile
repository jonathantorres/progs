release:
	rm -fr ./bin && mkdir ./bin
	GOOS=darwin GOARCH=amd64 go build && mv ./gonew ./bin/gonew_darwin
	GOOS=freebsd GOARCH=amd64 go build && mv ./gonew ./bin/gonew_freebsd
	GOOS=linux GOARCH=amd64 go build && mv ./gonew ./bin/gonew_linux
	GOOS=netbsd GOARCH=amd64 go build && mv ./gonew ./bin/gonew_netbsd
	GOOS=plan9 GOARCH=amd64 go build && mv ./gonew ./bin/gonew_plan9
	GOOS=windows GOARCH=amd64 go build && mv ./gonew.exe ./bin/gonew_windows.exe
