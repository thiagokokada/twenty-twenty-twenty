package core

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gioui.org/x/notify"
	ntf "github.com/thiagokokada/twenty-twenty-twenty/notification"
	"github.com/thiagokokada/twenty-twenty-twenty/sound"
)

var (
	Stop    context.CancelFunc
	loopCtx context.Context
	mu      sync.Mutex
)

type Settings struct {
	Duration  time.Duration
	Frequency time.Duration
	Pause     time.Duration
	Sound     bool
}

type Optional struct {
	Sound   bool
	Systray bool
}

func ParseFlags(
	progname string,
	args []string,
	version string,
	optional Optional,
) Settings {
	flags := flag.NewFlagSet(progname, flag.ExitOnError)
	durationInSec := flags.Uint(
		"duration",
		20,
		"how long each pause should be in seconds",
	)
	frequencyInSec := flags.Uint(
		"frequency",
		20*60,
		"how often the pause should be in seconds",
	)
	pauseInSec := new(uint)
	if optional.Systray {
		pauseInSec = flags.Uint(
			"pause",
			60*60,
			"how long the pause (from systray) should be in seconds",
		)
	}
	disableSound := new(bool)
	if optional.Sound {
		disableSound = flags.Bool(
			"disable-sound",
			false,
			"disable notification sound",
		)
	}
	showVersion := flags.Bool(
		"version",
		false,
		"print program version and exit",
	)
	flags.Parse(args)

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	return Settings{
		Duration:  time.Duration(*durationInSec) * time.Second,
		Frequency: time.Duration(*frequencyInSec) * time.Second,
		Pause:     time.Duration(*pauseInSec) * time.Second,
		Sound:     optional.Sound && !*disableSound,
	}
}

func Start(
	notifier notify.Notifier,
	settings *Settings,
) {
	mu.Lock()
	defer mu.Unlock()

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
	if loopCtx != nil {
		Stop() // make sure we cancel the previous instance
	}
	loopCtx, Stop = context.WithCancel(context.Background())
	loop(loopCtx, notifier, settings)
}

func Pause(
	ctx context.Context,
	notifier notify.Notifier,
	settings *Settings,
	timerCallback func(),
) {
	log.Printf("Pausing twenty-twenty-twenty for %.f hour...\n", settings.Pause.Hours())

	if loopCtx != nil {
		Stop() // cancelling current twenty-twenty-twenty goroutine
	}
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
		ntf.CancelAfter(cancelCtx, notification, &settings.Duration, &settings.Sound)
		go Start(notifier, settings)
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
				ntf.CancelAfter(cancelCtx, notification, &settings.Duration, &settings.Sound)
			}()
		case <-ctx.Done():
			log.Println("Disabling twenty-twenty-twenty...")
			cancelCtxCancel()
			return
		}
	}
}
