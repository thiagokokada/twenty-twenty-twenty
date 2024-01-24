os := $(shell uname -s)
arch := $(shell uname -m)

# Default to .app bundle in macOS since without a bundle the application
# doesn't work
.PHONY: all
ifeq ($(os),Darwin)
all: bin/TwentyTwentyTwenty_$(arch).app
else
all: bin/twenty-twenty-twenty
endif

.PHONY: run
ifeq ($(os),Darwin)
run: bin/TwentyTwentyTwenty_$(arch).app
		bin/TwentyTwentyTwenty_$(arch).app/Contents/MacOS/TwentyTwentyTwenty
else
run: bin/twenty-twenty-twenty
		bin/twenty-twenty-twenty
endif

.PHONY: lint
lint:
	test -z $(shell gofmt -l .)
	go vet -v ./...
	go run github.com/kisielk/errcheck -verbose ./...
	go run honnef.co/go/tools/cmd/staticcheck ./...

.PHONY: clean
clean:
	rm -rf bin

LDFLAGS := -X 'main.Version=$(shell git describe --tags --dirty)' -s -w

.PHONY: bin/twenty-twenty-twenty
bin/twenty-twenty-twenty:
	go build -v -ldflags="$(LDFLAGS)" -o $@

# Cross-build target for Windows:
# - bin/twenty-twenty-twenty-windows-386
# - bin/twenty-twenty-twenty-windows-arm64
# - bin/twenty-twenty-twenty-windows-amd64
.PHONY: bin/twenty-twenty-twenty-%.exe
bin/twenty-twenty-twenty-%.exe:
	GOOS=$(word 1,$(subst -, ,$*)) GOARCH=$(word 2,$(subst -, ,$*)) CGO_ENABLED=0 \
			 go build -v -ldflags="-H=windowsgui $(LDFLAGS)" -o $@

# Cross-build target, use as e.g.: `make bin/twenty-twenty-twenty-linux-arm64`
# Some valid targets:
# - bin/twenty-twenty-twenty-linux-amd64
# - bin/twenty-twenty-twenty-linux-arm64
# Since we set CGO_ENABLED=0, some features may be missing (e.g.: sound)
.PHONY: bin/twenty-twenty-twenty-%
bin/twenty-twenty-twenty-%:
	GOOS=$(word 1,$(subst -, ,$*)) GOARCH=$(word 2,$(subst -, ,$*)) CGO_ENABLED=0 \
			 go build -v -ldflags="$(LDFLAGS)" -o $@

# Nix target for static binaries in Linux
# Needs to be run from the same host that the binaries will be built
.PHONY: bin/twenty-twenty-twenty-linux-amd64-static
bin/twenty-twenty-twenty-linux-amd64-static:
	mkdir -p bin
	cp $(shell nix build '.#packages.x86_64-linux.twenty-twenty-twenty-static' --no-link --json | jq -r .[].outputs.out)/bin/twenty-twenty-twenty $@
	chmod +rwx $@

.PHONY: bin/twenty-twenty-twenty-linux-arm64-static
bin/twenty-twenty-twenty-linux-arm64-static:
	mkdir -p bin
	cp $(shell nix build '.#packages.aarch64-linux.twenty-twenty-twenty-static' --no-link --json | jq -r .[].outputs.out)/bin/twenty-twenty-twenty $@
	chmod +rwx $@

# macOS builds needs an .app bundle and (adhoc) signature to work
.PHONY: bin/TwentyTwentyTwenty_arm64.app
bin/TwentyTwentyTwenty_arm64.app:
	go run gioui.org/cmd/gogio -x -arch=arm64 -target=macos -ldflags="$(LDFLAGS)" -icon=./assets/eye.png -o=$@ .
	cp $@/Contents/Resources/icon.icns assets/macos/TwentyTwentyTwenty.app/Contents/Resources/icon.icns
	cp assets/macos/TwentyTwentyTwenty.app/Contents/Info.plist $@/Contents/Info.plist
	mv $@/Contents/MacOS/TwentyTwentyTwenty_arm64 $@/Contents/MacOS/TwentyTwentyTwenty
	codesign -s - $@

.PHONY: bin/TwentyTwentyTwenty_amd64.app
bin/TwentyTwentyTwenty_amd64.app:
	go run gioui.org/cmd/gogio -x -arch=amd64 -target=macos -ldflags="$(LDFLAGS)" -icon=./assets/eye.png -o=$@ .
	cp $@/Contents/Resources/icon.icns assets/macos/TwentyTwentyTwenty.app/Contents/Resources/icon.icns
	cp assets/macos/TwentyTwentyTwenty.app/Contents/Info.plist $@/Contents/Info.plist
	mv $@/Contents/MacOS/TwentyTwentyTwenty_amd64 $@/Contents/MacOS/TwentyTwentyTwenty
	codesign -s - $@

bin/TwentyTwentyTwenty_arm64.zip: bin/TwentyTwentyTwenty_arm64.app
	cd bin && zip -rv TwentyTwentyTwenty_arm64.zip TwentyTwentyTwenty_arm64.app

bin/TwentyTwentyTwenty_amd64.zip: bin/TwentyTwentyTwenty_amd64.app
	cd bin && zip -rv TwentyTwentyTwenty_amd64.zip TwentyTwentyTwenty_amd64.app
