all: bin/twenty-twenty-twenty-win-i386 bin/twenty-twenty-twenty-win-amd64 \
	bin/twenty-twenty-twenty-macos-arm64 bin/twenty-twenty-twenty-macos-amd64 \
	bin/twenty-twenty-twenty-linux-arm64 bin/twenty-twenty-twenty-linux-amd64

bin/twenty-twenty-twenty-win-i386: bin/twenty-twenty-twenty
	GOOS=windows GOARCH=386 go build -o $@

bin/twenty-twenty-twenty-win-amd64: bin/twenty-twenty-twenty
	GOOS=windows GOARCH=amd64 go build -o $@

bin/twenty-twenty-twenty-macos-arm64: bin/twenty-twenty-twenty
	GOOS=darwin GOARCH=arm64 go build -o $@

bin/twenty-twenty-twenty-macos-amd64: bin/twenty-twenty-twenty
	GOOS=darwin GOARCH=amd64 go build -o $@

bin/twenty-twenty-twenty-linux-arm64: bin/twenty-twenty-twenty
	GOOS=linux GOARCH=arm64 go build -o $@

bin/twenty-twenty-twenty-linux-amd64: bin/twenty-twenty-twenty
	GOOS=linux GOARCH=amd64 go build -o $@

bin/twenty-twenty-twenty: eye-solid.svg *.go go.mod go.sum
	go test
	go build -o $@
