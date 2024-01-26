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

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/vorbis"
)

const Enabled bool = true

// Maximum lag, good enough for this use case and will use lower CPU, but need
// to compesate the lag with time.Sleep() to not feel "strange" (e.g.: "floaty"
// notifications because the sound comes too late).
const lag time.Duration = time.Second / 4

var (
	buffer1     *beep.Buffer
	buffer2     *beep.Buffer
	mu          sync.Mutex
	wg          sync.WaitGroup
	initialised bool
	suspended   bool
	//go:embed assets/*.ogg
	notifications embed.FS
)

func PlaySendNotification(endCallback func()) {
	slog.Debug("Playing send notification sound")

	mu.Lock()
	wg.Add(1)
	err := speakerResume()
	mu.Unlock()

	if err != nil {
		log.Printf("Error while resuming speaker: %v\n", err)
		return
	}

	speaker.Play(beep.Seq(
		buffer1.Streamer(0, buffer1.Len()),
		beep.Callback(func() { wg.Done() }),
		beep.Callback(endCallback),
	))

	// compesate the lag
	time.Sleep(lag)
}

func PlayCancelNotification(endCallback func()) {
	slog.Debug("Playing cancel notification sound")

	mu.Lock()
	wg.Add(1)
	err := speakerResume()
	mu.Unlock()

	if err != nil {
		log.Printf("Error while resuming speaker: %v\n", err)
		return
	}

	speaker.Play(beep.Seq(
		buffer2.Streamer(0, buffer2.Len()),
		beep.Callback(func() { wg.Done() }),
		beep.Callback(endCallback),
	))

	// compesate the lag
	time.Sleep(lag)
}

func Init(suspend bool) (err error) {
	mu.Lock()
	defer mu.Unlock()

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

	slog.Debug(
		"Initialising speaker",
		"sampleRate", format.SampleRate,
		"bufferSize", format.SampleRate.N(lag),
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

	mu.Lock()
	defer mu.Unlock()

	// make sure that no sound is playing before suspending
	wg.Wait()
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

// WARN: this function is not thread safe, call it inside a mu.Lock()
func speakerResume() error {
	if !initialised {
		slog.Debug("Ignoring speaker resume call since it is not initialised yet")
		return nil
	}

	if suspended {
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

// WARN: this function is not thread safe, call it inside a mu.Lock()
func speakerSuspend() error {
	if !initialised {
		slog.Debug("Ignoring speaker suspend call since it is not initialised yet")
		return nil
	}

	if !suspended {
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
