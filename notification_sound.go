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

const notificationSoundEnabled bool = true

// Maximum lag, good enough for this use case and will use lower CPU, but need
// to compesate the lag with time.Sleep() to not feel "strange" (e.g.: "floaty"
// notifications because the sound comes too late).
// About as good we can get of CPU usage for now, until this issue is fixed:
// https://github.com/gopxl/beep/issues/137.
const lag time.Duration = time.Second

var (
	buffer1 *beep.Buffer
	buffer2 *beep.Buffer
	//go:embed assets/notification_1.ogg
	notification1 embed.FS
	//go:embed assets/notification_2.ogg
	notification2 embed.FS
)

func playSendNotificationSound() {
	done := make(chan bool)
	speaker.Play(
		beep.Seq(buffer1.Streamer(0, buffer1.Len())),
		beep.Callback(func() { done <- true }),
	)
	<-done
	// compesate the lag
	time.Sleep(lag)
}

func playCancelNotificationSound() {
	done := make(chan bool)
	speaker.Play(
		beep.Seq(buffer2.Streamer(0, buffer2.Len())),
		beep.Callback(func() { done <- true }),
	)
	<-done
	// compesate the lag
	time.Sleep(lag)
}

func initNotification() error {
	loadNotification := func(notification embed.FS, file string) (*beep.Buffer, beep.Format, error) {
		f, err := notification.Open(file)
		if err != nil {
			return nil, beep.Format{}, fmt.Errorf("load notification %s sound: %w", file, err)
		}
		streamer, format, err := vorbis.Decode(f)
		if err != nil {
			return nil, beep.Format{}, fmt.Errorf("decode notification %s sound: %w", file, err)
		}
		buffer := beep.NewBuffer(format)
		buffer.Append(streamer)
		return buffer, format, nil
	}

	var format beep.Format
	var err error

	buffer1, format, err = loadNotification(notification1, "assets/notification_1.ogg")
	if err != nil {
		return fmt.Errorf("notification 1 sound failed: %w", err)
	}

	// ignoring format since all audio files should have the same format
	buffer2, _, err = loadNotification(notification2, "assets/notification_2.ogg")
	if err != nil {
		return fmt.Errorf("notification 2 sound failed: %w", err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(lag))

	return nil
}
