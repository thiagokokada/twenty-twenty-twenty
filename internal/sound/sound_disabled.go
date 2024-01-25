//go:build !windows && !darwin && !cgo
// +build !windows,!darwin,!cgo

package sound

import "time"

const Enabled bool = false
const panicMsg string = "Sound disabled in this build"

func Resume() error {
	panic(panicMsg)
}

func Suspend() error {
	panic(panicMsg)
}

func SuspendAfter(time.Duration) error {
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
