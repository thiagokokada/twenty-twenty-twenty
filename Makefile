.PHONY: all

os := $(shell uname -s)
arch := $(shell uname -a)

ifeq ($(os),Darwin)
ifeq ($(arch),arm64)
all: bin/TwentyTwentyTwenty_arm64.app bin/TwentyTwentyTwenty_amd64.app
else
all: bin/TwentyTwentyTwenty_amd64.app
endif
else
all: bin/twenty-twenty-twenty
endif

# Cross-build target for Windows:
# - bin/twenty-twenty-twenty-windows-386
# - bin/twenty-twenty-twenty-windows-arm64
# - bin/twenty-twenty-twenty-windows-amd64
bin/twenty-twenty-twenty-%.exe: assets/* *.go go.mod go.sum
	GOOS=$(word 1,$(subst -, ,$*)) GOARCH=$(word 2,$(subst -, ,$*)) CGO_ENABLED=0 \
			 go build -v -ldflags="-H=windowsgui -X 'main.Version=$(shell git describe --tags --dirty)' -s -w" -o $@

# Cross-build target, use as e.g.: `make bin/twenty-twenty-twenty-linux-arm64`
# Some valid targets:
# - bin/twenty-twenty-twenty-linux-amd64 # no audio
# - bin/twenty-twenty-twenty-linux-arm64 # no audio
# - bin/twenty-twenty-twenty-freebsd-amd64 # no audio
bin/twenty-twenty-twenty-%: assets/* *.go go.mod go.sum
	GOOS=$(word 1,$(subst -, ,$*)) GOARCH=$(word 2,$(subst -, ,$*)) CGO_ENABLED=0 \
			 go build -v -ldflags="-X 'main.Version=$(shell git describe --tags --dirty)' -s -w" -o $@

# Not including the `-s -w` flags here since they're important for debugging
# and this target is mostly used for development
bin/twenty-twenty-twenty: assets/* *.go go.mod go.sum
	go build -v -ldflags="-X 'main.Version=$(shell git describe --tags --dirty)'" -o $@

bin/TwentyTwentyTwenty_arm64.app: assets/* *.go go.mod go.sum
	go generate loop_darwin.go
	mkdir -p bin/
	rm -rf bin/TwentyTwentyTwenty_*.app
	mv TwentyTwentyTwenty.app/TwentyTwentyTwenty_*.app bin/
	rmdir TwentyTwentyTwenty.app
	cp bin/TwentyTwentyTwenty_arm64.app/Contents/Resources/icon.icns assets/macos/TwentyTwentyTwenty.app/Contents/Resources/icon.icns


bin/TwentyTwentyTwenty_amd64.app: bin/TwentyTwentyTwenty_arm64.app

bin/twenty-twenty-twenty-linux-amd64-static: assets/* *.go go.mod go.sum *.nix
	cp $(shell nix build '.#packages.x86_64-linux.twenty-twenty-twenty-static' --no-link --json | jq -r .[].outputs.out)/bin/twenty-twenty-twenty $@

bin/twenty-twenty-twenty-linux-arm64-static: assets/* *.go go.mod go.sum *.nix
	cp $(shell nix build '.#packages.aarch64-linux.twenty-twenty-twenty-static' --no-link --json | jq -r .[].outputs.out)/bin/twenty-twenty-twenty $@

clean:
	rm -rf bin
