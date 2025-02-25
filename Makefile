os := $(shell uname -s)
arch := $(shell uname -m)

# Default to .app bundle in macOS since without a bundle the application
# doesn't work
.PHONY: all
ifeq ($(os)-$(arch),Darwin-arm64)
all: bin/TwentyTwentyTwenty_arm64.app
else ifeq ($(os)-$(arch),Darwin-x86_64)
all: bin/TwentyTwentyTwenty_amd64.app
else
all: bin/twenty-twenty-twenty
endif

.PHONY: run
ifeq ($(os)-$(arch),Darwin-arm64)
run: bin/TwentyTwentyTwenty_arm64.app
		bin/TwentyTwentyTwenty_arm64.app/Contents/MacOS/TwentyTwentyTwenty
else ifeq ($(os)-$(arch),Darwin-x86_64)
run: bin/TwentyTwentyTwenty_amd64.app
		bin/TwentyTwentyTwenty_amd64.app/Contents/MacOS/TwentyTwentyTwenty
else
run: bin/twenty-twenty-twenty
		bin/twenty-twenty-twenty
endif

.PHONY: coverage
coverage:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html

.PHONY: lint
lint:
	test -z $(shell gofmt -l .)
	go vet -v ./...
	go tool github.com/kisielk/errcheck -verbose ./...
	go tool honnef.co/go/tools/cmd/staticcheck ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: test-ci
test-ci:
	CI=1 go test -race -v ./...

.PHONY: clean
clean:
	rm -rf bin cover.*

LDFLAGS := -X 'main.Version=$(shell git describe --tags --dirty)' -s -w

.PHONY: bin/twenty-twenty-twenty
bin/twenty-twenty-twenty:
	go build -v -ldflags="$(LDFLAGS)" -o $@

# Cross-build target for Windows
bin/twenty-twenty-twenty-windows-%.exe: PHONY_TARGET
	GOOS=windows GOARCH=$* CGO_ENABLED=0 go build -v -ldflags="-H=windowsgui $(LDFLAGS)" -o $@

bin/twenty-twenty-twenty-windows-386.exe:
bin/twenty-twenty-twenty-windows-amd64.exe:
bin/twenty-twenty-twenty-windows-arm64.exe:

# Cross-build target, use as e.g.: `make bin/twenty-twenty-twenty-linux-arm64`
# Since we set CGO_ENABLED=0, some features may be missing (e.g.: sound), but
# the binaries are static
bin/twenty-twenty-twenty-%: PHONY_TARGET
	GOOS=$(word 1,$(subst -, ,$*)) GOARCH=$(word 2,$(subst -, ,$*)) CGO_ENABLED=0 go build -v -ldflags="$(LDFLAGS)" -o $@

bin/twenty-twenty-twenty-linux-amd64:
bin/twenty-twenty-twenty-linux-arm64:

# Nix target for static binaries in Linux with CGO_ENABLED
# Needs to have nix installed and to be run from the same host that the
# binaries will be built
.PHONY: bin/twenty-twenty-twenty-linux-amd64-static
bin/twenty-twenty-twenty-linux-amd64-static:
	mkdir -p bin
	cp $(shell nix build '.#packages.x86_64-linux.twenty-twenty-twenty-static' --no-link --json | jq -r .[].outputs.out)/bin/twenty-twenty-twenty $@
	chmod +rwx $@

.PHONY: bin/twenty-twenty-twenty-linux-arm64-static
bin/twenty-twenty-twenty-linux-arm64-static:
	mkdir -p bin
	cp $(shell nix build '.#packages.aarch64-linux.twenty-twenty-twenty-static' --no-link --json | jq -r .[].outputs.out)/bin/twenty-twenty-twenty $@

# macOS builds needs an .app bundle and (adhoc) signature to work, and only
# work in macOS itself (since it needs CGO_ENABLED and codesign)
bin/TwentyTwentyTwenty_%.app: PHONY_TARGET
	go tool gioui.org/cmd/gogio -x -arch=$* -target=macos -ldflags="$(LDFLAGS)" -icon=./assets/eye.png -o=$@ .
	cp $@/Contents/Resources/icon.icns assets/macos/TwentyTwentyTwenty.app/Contents/Resources/icon.icns
	cp assets/macos/TwentyTwentyTwenty.app/Contents/Info.plist $@/Contents/Info.plist
	mv $@/Contents/MacOS/TwentyTwentyTwenty_$* $@/Contents/MacOS/TwentyTwentyTwenty
	codesign -s - $@

bin/TwentyTwentyTwenty_arm64.app:
bin/TwentyTwentyTwenty_amd64.app:

bin/TwentyTwentyTwenty_%.zip: bin/TwentyTwentyTwenty_%.app
	cd bin && zip -rv TwentyTwentyTwenty_$*.zip TwentyTwentyTwenty_$*.app

bin/TwentyTwentyTwenty_arm64.zip:
bin/TwentyTwentyTwenty_amd64.zip:

# To be used for targets with pattern (e.g.: %) since Makefile doesn't
# understand patterns in PHONY targets
.PHONY: PHONY_TARGET
PHONY_TARGET:
