//go:build windows || darwin || cgo
// +build windows darwin cgo

package sound

import (
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/vorbis"
)

const Enabled bool = true

// Maximum sound notification lag, 1000ms / 10 = 100ms
const lag time.Duration = time.Second / 10

var (
	buffer1 *beep.Buffer
	buffer2 *beep.Buffer
	//go:embed assets/*.ogg
	notifications embed.FS
	initialized   bool
)

func speakerResume() {
	err := speaker.Resume()
	if err != nil {
		log.Printf("Error while resuming speaker: %v\n", err)
	}
}

func speakerSuspend() {
	err := speaker.Suspend()
	if err != nil {
		log.Printf("Error while suspending speaker: %v\n", err)
	}
}

func PlaySendNotification(callback func()) {
	speakerResume()

	speaker.Play(beep.Seq(
		buffer1.Streamer(0, buffer1.Len()),
		beep.Callback(speakerSuspend),
		beep.Callback(callback),
	))
}

func PlayCancelNotification(callback func()) {
	speakerResume()

	speaker.Play(beep.Seq(
		buffer2.Streamer(0, buffer2.Len()),
		beep.Callback(speakerSuspend),
		beep.Callback(callback),
	))
}

func Init() (err error) {
	// should be safe to call multiple times
	if !initialized {
		var format beep.Format

		buffer1, format, err = loadSound("assets/notification_1.ogg")
		if err != nil {
			return fmt.Errorf("notification 1 sound failed: %w", err)
		}
		// ignoring format since all audio files should have the same format
		buffer2, _, err = loadSound("assets/notification_2.ogg")
		if err != nil {
			return fmt.Errorf("notification 2 sound failed: %w", err)
		}

		err = speaker.Init(format.SampleRate, format.SampleRate.N(lag))
		if err != nil {
			return fmt.Errorf("speaker init: %w", err)
		}
		err = speaker.Suspend()
		if err != nil {
			return fmt.Errorf("speaker suspend: %w", err)
		}
		initialized = true

		return nil
	}
	return nil
}

func loadSound(file string) (*beep.Buffer, beep.Format, error) {
	f, err := notifications.Open(file)
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
