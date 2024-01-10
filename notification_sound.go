//go:build windows || darwin || cgo
// +build windows darwin cgo

package main

import (
	"embed"
	"log"
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

func playNotificationSound() {
	done := make(chan bool)
	speaker.Play(
		beep.Seq(Buffer.Streamer(0, Buffer.Len())),
		beep.Callback(func() { done <- true }),
	)
	<-done
}

func initBeep() {
	f, err := NotificationSound.Open("notification.ogg")
	if err != nil {
		log.Fatalf("Failed to load notification sound: %v\n", err)
	}

	streamer, format, err := vorbis.Decode(f)
	if err != nil {
		log.Fatalf("Failed to decode the notification sound: %v\n", err)
	}
	Buffer = beep.NewBuffer(format)
	Buffer.Append(streamer)

	// 1s/4 = 250ms of lag, good enough for this use case
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/4))
}