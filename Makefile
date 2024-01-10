.PHONY: all

ifeq ($(OS),Windows_NT)
    detected_OS := Windows
else
    detected_OS := $(shell uname -s)
endif

ifeq ($(detected_OS),Darwin)
all: bin/TwentyTwentyTwenty_amd64.app bin/TwentyTwentyTwenty_amd64.app
else
all: bin/twenty-twenty-twenty
endif

bin/twenty-twenty-twenty: *.go go.mod go.sum
	 go build -v -ldflags="-X 'main.Version=$(shell git describe --tags --dirty)'" -o $@

bin/TwentyTwentyTwenty_amd64.app: bin/TwentyTwentyTwenty_aarch64.app

bin/TwentyTwentyTwenty_aarch64.app: eye.png *.go go.mod go.sum
	go generate loop_darwin.go
	mkdir -p bin/
	mv TwentyTwentyTwenty.app/TwentyTwentyTwenty_*.app bin/
	rm -rf TwentyTwentyTwenty.app

clean:
	rm -rf bin
