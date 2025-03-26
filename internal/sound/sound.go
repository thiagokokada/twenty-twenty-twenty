//go:build !nosound && cgo
// +build !nosound,cgo

package sound

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
)

const Enabled bool = true

// Maximum lag, good enough for this use case and will use lower CPU, but need
// to compesate the lag with time.Sleep() to not feel "strange" (e.g.: "floaty"
// notifications because the sound comes too late).
const lag time.Duration = time.Second / 4

var (
	mu        sync.Mutex
	suspended bool

	initialised bool
	sound1      sound
	sound2      sound
	wg          wgCount
	//go:embed assets/*.ogg
	notifications embed.FS
)

func PlaySendNotification(endCallback func()) {
	slog.Debug("Playing send notification sound", "wg", wg.count())

	err := playSound(sound1, endCallback)
	if err != nil {
		log.Printf("Error while playing send notification sound: %v\n", err)
	}

	// compesate the lag
	time.Sleep(lag)
}

func PlayCancelNotification(endCallback func()) {
	slog.Debug("Playing cancel notification sound", "wg", wg.count())

	err := playSound(sound2, endCallback)
	if err != nil {
		log.Printf("Error while playing cancel notification sound: %v\n", err)
	}

	// compesate the lag
	time.Sleep(lag)
}

func Init(suspend bool) error {
	if initialised {
		slog.Debug("Sound already initialised")
		return nil
	}

	buffer, format, err := loadSound("assets/notification_1.ogg")
	if err != nil {
		return fmt.Errorf("notification 1 sound failed: %w", err)
	}
	sound1 = sound{buffer: buffer, name: "send"}

	// ignoring format since all audio files should have the same format
	buffer, _, err = loadSound("assets/notification_2.ogg")
	if err != nil {
		return fmt.Errorf("notification 2 sound failed: %w", err)
	}
	sound2 = sound{buffer: buffer, name: "cancel"}

	slog.Debug(
		"Initialising speaker",
		"sampleRate", format.SampleRate,
		"bufferSize", format.SampleRate.N(lag),
		"lag", lag,
	)
	err = speaker.Init(format.SampleRate, format.SampleRate.N(lag))
	if err != nil {
		return fmt.Errorf("speaker init: %w", err)
	}
	initialised = true

	if suspend {
		err = speakerSuspend()
		if err != nil {
			return fmt.Errorf("speaker suspend: %w", err)
		}
	}

	return nil
}

func SuspendAfter(after time.Duration) {
	slog.Debug("Suspending sound", "afterSeconds", after.Seconds())
	time.Sleep(after)

	slog.Debug("Waiting sounds to finish playing before suspending", "wg", wg.count())
	wg.Wait()
	slog.Debug("Finished playing sound, calling speaker suspend", "wg", wg.count())

	err := speakerSuspend()
	if err != nil {
		log.Printf("Error while suspending speaker: %v\n", err)
	}
}

func loadSound(file string) (*beep.Buffer, beep.Format, error) {
	slog.Debug("Loading sound", "file", file)

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

func playSound(s sound, endCallback func()) error {
	wg.add(1)

	err := speakerResume()
	if err != nil {
		return fmt.Errorf("play sound resume: %w", err)
	}

	speaker.Play(beep.Seq(
		s.buffer.Streamer(0, s.buffer.Len()),
		beep.Callback(func() {
			wg.done()
			slog.Debug("Notification sound done", "sound", s, "wg", wg.count())
		}),
		beep.Callback(endCallback),
	))
	return nil
}

func speakerResume() error {
	if !initialised {
		slog.Debug("Ignoring speaker resume call since it is not initialised yet")
		return nil
	}

	if suspended {
		mu.Lock()
		defer mu.Unlock()

		slog.Debug("Resuming speaker")
		err := speaker.Resume()
		if err != nil {
			return fmt.Errorf("resuming speaker: %w", err)
		}
		suspended = false
	} else {
		slog.Debug("Speaker already resumed")
	}
	return nil
}

func speakerSuspend() error {
	if !initialised {
		slog.Debug("Ignoring speaker suspend call since it is not initialised yet")
		return nil
	}

	if !suspended {
		mu.Lock()
		defer mu.Unlock()

		slog.Debug("Suspending speaker to reduce CPU usage")
		speaker.Clear()
		err := speaker.Suspend()
		if err != nil {
			return fmt.Errorf("suspending speaker: %w", err)
		}
		suspended = true
	} else {
		slog.Debug("Speaker already suspended")
	}
	return nil
}
