//go:build tools

//go:generate sh -c "go install $(go list -e -f '{{join .Imports \" \"}}' tools.go)"

package main

import _ "gioui.org/cmd/gogio"
