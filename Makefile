.PHONY: all lint clean

os := $(shell uname -s)
arch := $(shell uname -m)

ifeq ($(os),Darwin)
ifeq ($(arch),arm64)
all: bin/TwentyTwentyTwenty_arm64.app
else
all: bin/TwentyTwentyTwenty_amd64.app
endif
else
all: bin/twenty-twenty-twenty
endif

LDFLAGS := -X 'main.Version=$(shell git describe --tags --dirty)' -s -w
# icon.icns is always updated so ignore it from dependencies
DEPS := $(shell find assets/* -type f ! -name icon.icns) *.go go.mod go.sum

# Cross-build target for Windows:
# - bin/twenty-twenty-twenty-windows-386
# - bin/twenty-twenty-twenty-windows-arm64
# - bin/twenty-twenty-twenty-windows-amd64
bin/twenty-twenty-twenty-%.exe: $(DEPS)
	GOOS=$(word 1,$(subst -, ,$*)) GOARCH=$(word 2,$(subst -, ,$*)) CGO_ENABLED=0 \
			 go build -v -ldflags="-H=windowsgui $(LDFLAGS)" -o $@

# Cross-build target, use as e.g.: `make bin/twenty-twenty-twenty-linux-arm64`
# Some valid targets:
# - bin/twenty-twenty-twenty-linux-amd64
# - bin/twenty-twenty-twenty-linux-arm64
# Since we set CGO_ENABLED=0, some features may be missing (e.g.: sound)
bin/twenty-twenty-twenty-%: $(DEPS)
	GOOS=$(word 1,$(subst -, ,$*)) GOARCH=$(word 2,$(subst -, ,$*)) CGO_ENABLED=0 \
			 go build -v -ldflags="$(LDFLAGS)" -o $@

bin/twenty-twenty-twenty: assets/* *.go go.mod go.sum
	go build -v -ldflags="$(LDFLAGS)" -o $@

bin/TwentyTwentyTwenty_arm64.zip: bin/TwentyTwentyTwenty_arm64.app
	cd bin && zip -rv TwentyTwentyTwenty_arm64.zip TwentyTwentyTwenty_arm64.app

bin/TwentyTwentyTwenty_amd64.zip: bin/TwentyTwentyTwenty_amd64.app
	cd bin && zip -rv TwentyTwentyTwenty_amd64.zip TwentyTwentyTwenty_amd64.app

bin/TwentyTwentyTwenty_arm64.app: $(DEPS)
	go run gioui.org/cmd/gogio -arch=arm64 -target=macos -ldflags="$(LDFLAGS)" -icon=./assets/eye.png -o=$@ .
	cp assets/macos/TwentyTwentyTwenty.app/Contents/Info.plist $@/Contents/Info.plist
	mv $@/Contents/MacOS/TwentyTwentyTwenty_arm64 $@/Contents/MacOS/TwentyTwentyTwenty
	codesign -s - $@

bin/TwentyTwentyTwenty_amd64.app: $(DEPS)
	go run gioui.org/cmd/gogio -arch=amd64 -target=macos -ldflags="$(LDFLAGS)" -icon=./assets/eye.png -o=$@ .
	cp assets/macos/TwentyTwentyTwenty.app/Contents/Info.plist $@/Contents/Info.plist
	mv $@/Contents/MacOS/TwentyTwentyTwenty_amd64 $@/Contents/MacOS/TwentyTwentyTwenty
	codesign -s - $@

bin/twenty-twenty-twenty-linux-amd64-static: $(DEPS) *.nix
	cp $(shell nix build '.#packages.x86_64-linux.twenty-twenty-twenty-static' --no-link --json | jq -r .[].outputs.out)/bin/twenty-twenty-twenty $@

bin/twenty-twenty-twenty-linux-arm64-static: $(DEPS) *.nix
	cp $(shell nix build '.#packages.aarch64-linux.twenty-twenty-twenty-static' --no-link --json | jq -r .[].outputs.out)/bin/twenty-twenty-twenty $@

lint:
	go vet -v ./...
	go run github.com/kisielk/errcheck -verbose ./...
	go run honnef.co/go/tools/cmd/staticcheck ./...

clean:
	rm -rf bin
