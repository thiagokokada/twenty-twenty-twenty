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

# Cross-build target, use as e.g.: `make bin/twenty-twenty-twenty-linux-arm64`
# Some valid targets:
# - bin/twenty-twenty-twenty-windows-386
# - bin/twenty-twenty-twenty-windows-amd64
# - bin/twenty-twenty-twenty-linux-amd64 # no audio
# - bin/twenty-twenty-twenty-linux-arm64 # no audio
# - bin/twenty-twenty-twenty-freebsd-amd64 # no audio
bin/twenty-twenty-twenty-%: *.go go.mod go.sum
	GOOS=$(word 1,$(subst -, ,$*)) GOARCH=$(word 2,$(subst -, ,$*)) CGO_ENABLED=0 \
	     go build -v -ldflags="-X 'main.Version=$(shell git describe --tags --dirty)'" -o $@

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
