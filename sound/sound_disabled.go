//go:build !windows && !darwin && !cgo
// +build !windows,!darwin,!cgo

package sound

const Enabled bool = false

func PlaySendNotification() {
	panic("Not implemented")
}

func PlayCancelNotification() {
	panic("Not implemented")
}

func Init() error {
	panic("Not implemented")
}
