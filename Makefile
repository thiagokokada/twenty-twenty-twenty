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

bin/twenty-twenty-twenty: assets/* *.go go.mod go.sum
	 go build -v -ldflags="-X 'main.Version=$(shell git describe --tags --dirty)'" -o $@

bin/TwentyTwentyTwenty_arm64.app: assets/* *.go go.mod go.sum
	go generate loop_darwin.go
	mkdir -p bin/
	rm -rf bin/TwentyTwentyTwenty_*.app
	mv TwentyTwentyTwenty.app/TwentyTwentyTwenty_*.app bin/
	rmdir TwentyTwentyTwenty.app

bin/TwentyTwentyTwenty_amd64.app: bin/TwentyTwentyTwenty_arm64.app

clean:
	rm -rf bin
