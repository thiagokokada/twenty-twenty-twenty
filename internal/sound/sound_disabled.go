//go:build !windows && !darwin && !cgo
// +build !windows,!darwin,!cgo

package sound

import "time"

const Enabled bool = false
const panicMsg string = "Sound disabled in this build"

func Resume() {
	panic(panicMsg)
}

func Suspend() {
	panic(panicMsg)
}

func SuspendAfter(time.Duration) {
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
