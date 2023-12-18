all: bin/twenty-twenty-twenty-windows-386 bin/twenty-twenty-twenty-windows-amd64 \
	bin/twenty-twenty-twenty-darwin-arm64 bin/twenty-twenty-twenty-darwin-amd64 \
	bin/twenty-twenty-twenty-linux-arm64 bin/twenty-twenty-twenty-linux-amd64

bin/twenty-twenty-twenty-%: bin/twenty-twenty-twenty
	GOOS=$(word 1,$(subst -, ,$*)) GOARCH=$(word 2,$(subst -, ,$*)) go build -o $@

bin/twenty-twenty-twenty: eye-solid.svg *.go go.mod go.sum
	go test
	go build -o $@
