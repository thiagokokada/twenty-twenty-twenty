//go:build nosound || !cgo
// +build nosound !cgo

package sound

import "time"

const Enabled bool = false
const panicMsg string = "Sound disabled in this build"

func SuspendAfter(after time.Duration) {
	panic(panicMsg)
}

func PlaySendNotification(func()) {
	panic(panicMsg)
}

func PlayCancelNotification(func()) {
	panic(panicMsg)
}

func Init(bool) error {
	panic(panicMsg)
}
