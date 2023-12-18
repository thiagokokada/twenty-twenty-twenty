.PHONY: clean
export CGO_ENABLED := 1

all: bin/twenty-twenty-twenty

TwentyTwentyTwenty.app: bin/twenty-twenty-twenty
	gogio -target macos -icon ./eye.png -o $@ .
	codesign -s - $@/*.app

bin/twenty-twenty-twenty: eye.png *.go go.mod go.sum
	go build -o $@

clean:
	rm -rf bin TwentyTwentyTwenty.app
