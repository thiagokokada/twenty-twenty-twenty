//go:build windows || darwin || cgo
// +build windows darwin cgo

package main

import (
	"embed"
	"fmt"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/vorbis"
)

var (
	//go:embed notification.ogg
	NotificationSound        embed.FS
	Buffer                   *beep.Buffer
	notificationSoundEnabled = true
)

func playNotificationSound() chan bool {
	done := make(chan bool)
	speaker.Play(
		beep.Seq(Buffer.Streamer(0, Buffer.Len())),
		beep.Callback(func() { done <- true }),
	)
	return done
}

func initBeep() error {
	f, err := NotificationSound.Open("notification.ogg")
	if err != nil {
		return fmt.Errorf("failed to load notification sound: %w", err)
	}

	streamer, format, err := vorbis.Decode(f)
	if err != nil {
		return fmt.Errorf("failed to decode the notification sound: %w", err)
	}
	Buffer = beep.NewBuffer(format)
	Buffer.Append(streamer)

	// 1s/4 = 250ms of lag, good enough for this use case
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/4))

	return nil
}
