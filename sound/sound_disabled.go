//go:build !windows && !darwin && !cgo
// +build !windows,!darwin,!cgo

package sound

const Enabled bool = false
const panicMsg string = "Sound disabled in this build"

func PlaySendNotification() {
	panic(panicMsg)
}

func PlayCancelNotification() {
	panic(panicMsg)
}

func Init() error {
	panic(panicMsg)
}
