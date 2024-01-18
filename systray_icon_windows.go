//go:build !nosystray
// +build !nosystray

package main

import _ "embed"

//go:embed assets/eye_light.ico
var systrayIcon []byte
