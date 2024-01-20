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
)

var (
	Stop context.CancelFunc
	Ctx  context.Context
	mu   sync.Mutex
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
	optional Optional,
) {
	mu.Lock()
	defer mu.Unlock()

	if optional.Sound {
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
	if Ctx != nil {
		Stop() // make sure we cancel the previous instance
	}
	Ctx, Stop = context.WithCancel(context.Background())
	go loop(Ctx, notifier, settings)
}

func Pause(
	ctx context.Context,
	notifier notify.Notifier,
	settings *Settings,
	optional Optional,
	timerCallback func(),
) {
	log.Printf("Pausing twenty-twenty-twenty for %.f hour...\n", settings.Pause.Hours())

	if Ctx != nil {
		Stop() // cancelling current twenty-twenty-twenty goroutine
	}
	timer := time.NewTimer(settings.Pause)

	select {
	case <-timer.C:
		err := ntf.SendWithDuration(
			ctx,
			notifier,
			&settings.Duration,
			&settings.Sound,
			"Resuming 20-20-20",
			fmt.Sprintf("You will see a notification every %.f minutes(s)", settings.Frequency.Minutes()),
		)
		if err != nil {
			log.Fatalf("Error while resuming notification: %v. Exiting...\n", err)
		}
		Start(notifier, settings, optional)
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
			log.Println("Sending notification...")
			err := ntf.SendWithDuration(
				cancelCtx,
				notifier,
				&settings.Duration,
				&settings.Sound,
				"Time to rest your eyes",
				fmt.Sprintf("Look at 20 feet (~6 meters) away for %.f seconds", settings.Duration.Seconds()),
			)
			if err != nil {
				log.Printf("Error while sending notification: %v.\n", err)
			}
		case <-ctx.Done():
			log.Println("Disabling twenty-twenty-twenty...")
			cancelCtxCancel()
			return
		}
	}
}
