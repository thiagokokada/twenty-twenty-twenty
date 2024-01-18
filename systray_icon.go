//go:build !nosystray && !windows
// +build !nosystray,!windows

package main

import _ "embed"

//go:embed assets/eye_light.png
var systrayIcon []byte
