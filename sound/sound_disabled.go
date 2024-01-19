//go:build !windows && !darwin && !cgo
// +build !windows,!darwin,!cgo

package sound

const Enabled bool = false
const panicMsg string = "Sound disabled in this build"

func PlaySendNotification(callback func()) {
	panic(panicMsg)
}

func PlayCancelNotification(callback func()) {
	panic(panicMsg)
}

func Init() error {
	panic(panicMsg)
}
