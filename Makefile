export CGO_ENABLED := 0

all: bin/twenty-twenty-twenty-windows-386 bin/twenty-twenty-twenty-windows-amd64 \
	bin/twenty-twenty-twenty-linux-arm64 bin/twenty-twenty-twenty-linux-amd64

bin/twenty-twenty-twenty-%: *.go go.mod go.sum
	GOOS=$(word 1,$(subst -, ,$*)) GOARCH=$(word 2,$(subst -, ,$*)) go build -o $@

TwentyTwentyTwenty.app: eye.png *.go go.mod go.sum
	CGO_ENABLED=1 go generate loop_darwin.go

clean:
	rm -rf bin TwentyTwentyTwenty.app
