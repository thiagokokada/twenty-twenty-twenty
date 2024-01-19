package core

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gioui.org/x/notify"
	ntf "github.com/thiagokokada/twenty-twenty-twenty/notification"
	"github.com/thiagokokada/twenty-twenty-twenty/sound"
)

var (
	ctx  context.Context
	Stop context.CancelFunc
)

type Settings struct {
	Duration  time.Duration
	Frequency time.Duration
	Pause     time.Duration
	Sound     bool
}

func ParseFlags(
	version string,
	systrayEnabled bool,
	soundEnabled bool,
) Settings {
	durationInSec := flag.Uint(
		"duration",
		20,
		"how long each pause should be in seconds",
	)
	frequencyInSec := flag.Uint(
		"frequency",
		20*60,
		"how often the pause should be in seconds",
	)
	pauseInSec := new(uint)
	if systrayEnabled {
		pauseInSec = flag.Uint(
			"pause",
			60*60,
			"how long the pause (from systray) should be in seconds",
		)
	}
	disableSound := new(bool)
	if soundEnabled {
		disableSound = flag.Bool(
			"disable-sound",
			false,
			"disable notification sound",
		)
	}
	showVersion := flag.Bool(
		"version",
		false,
		"print program version and exit",
	)
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	return Settings{
		Duration:  time.Duration(*durationInSec) * time.Second,
		Frequency: time.Duration(*frequencyInSec) * time.Second,
		Pause:     time.Duration(*pauseInSec) * time.Second,
		Sound:     soundEnabled && !*disableSound,
	}
}

func Start(
	notifier notify.Notifier,
	settings *Settings,
) {
	if sound.Enabled {
		log.Printf(
			"Running twenty-twenty-twenty every %.1f minute(s), with %.f second(s) duration and sound set to %t...\n",
			settings.Frequency.Minutes(),
			settings.Duration.Seconds(),
			settings.Sound,
		)
	} else {
		log.Printf(
			"Running twenty-twenty-twenty every %.1f minute(s), with %.f second(s) duration...\n",
			settings.Frequency.Minutes(),
			settings.Duration.Seconds(),
		)
	}

	if ctx != nil {
		Stop() // make sure we cancel the previous instance
	}
	ctx, Stop = context.WithCancel(context.Background())
	go loop(ctx, notifier, settings)
}

func Pause(
	ctx context.Context,
	notifier notify.Notifier,
	settings *Settings,
	timerCallback func(),
) {
	log.Printf("Pausing twenty-twenty-twenty for %.f hour...\n", settings.Pause.Hours())
	Stop() // cancelling current twenty-twenty-twenty goroutine
	timer := time.NewTimer(settings.Pause)
	// context to the resuming notification cancellation, since the program
	// may be paused or disabled again before the notification finishes
	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	select {
	case <-timer.C:
		notification := ntf.Send(
			notifier,
			"Resuming 20-20-20",
			fmt.Sprintf("You will see a notification every %.f minutes(s)", settings.Frequency.Minutes()),
			&settings.Sound,
		)
		if notification == nil {
			log.Printf("Resume notification failed...")
		}
		go ntf.CancelAfter(cancelCtx, notification, &settings.Duration, &settings.Sound)
		Start(notifier, settings)
		timerCallback()
	case <-ctx.Done():
	}
}

func loop(
	ctx context.Context,
	notifier notify.Notifier,
	settings *Settings,
) {
	ticker := time.NewTicker(settings.Frequency)
	cancelCtx, cancelCtxCancel := context.WithCancel(context.Background())
	for {
		select {
		case <-ticker.C:
			go func() {
				log.Println("Sending notification...")
				notification := ntf.Send(
					notifier,
					"Time to rest your eyes",
					fmt.Sprintf("Look at 20 feet (~6 meters) away for %.f seconds", settings.Duration.Seconds()),
					&settings.Sound,
				)
				go ntf.CancelAfter(cancelCtx, notification, &settings.Duration, &settings.Sound)
			}()
		case <-ctx.Done():
			log.Println("Disabling twenty-twenty-twenty...")
			cancelCtxCancel()
			return
		}
	}
}
