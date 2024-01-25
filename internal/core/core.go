package core

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/thiagokokada/twenty-twenty-twenty/internal/notification"
	"github.com/thiagokokada/twenty-twenty-twenty/internal/sound"
)

var (
	cancelLoopCtx context.CancelFunc
	loopCtx       context.Context
	mu            sync.Mutex
)

/*
Settings struct.

'Duration' will be the duration of each notification. For example, if is 20
seconds, it means that each notification will stay by 20 seconds.

'Frequency' is how often each notification will be shown. For example, if it is
20 minutes, a new notification will appear at every 20 minutes.

'Pause' is the duration of the pause. For example, if it is 1 hour, we will
disable notifications for 1 hour.

'Sound' enables or disables sound every time a notification is shown.
*/
type Settings struct {
	Duration  time.Duration
	Frequency time.Duration
	Pause     time.Duration
	Sound     bool
}

/*
Optional struct.

This is used for features that are optional in the program, for example if sound
or systray are permanently disabled.
*/
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
	// using flag.ExitOnError, so no error will ever be returned
	_ = flags.Parse(args)

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

/*
Returns the current running twenty-twenty-twenty's context.

Can be used to register for context cancellation, so if [Stop] is called the
context will be done (see [pkg/context] for details).
*/
func Ctx() context.Context { return loopCtx }

/*
Start twenty-twenty-twenty.

This will start the main twenty-twenty-twenty loop in a goroutine, so avoid
calling this function inside a goroutine.
*/
func Start(
	settings *Settings,
	optional Optional,
) {
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
	Stop() // make sure we cancel the previous instance

	mu.Lock()
	defer mu.Unlock()
	loopCtx, cancelLoopCtx = context.WithCancel(context.Background())
	go loop(loopCtx, settings)
}

/*
Stop twenty-twenty-twenty.
*/
func Stop() {
	mu.Lock()
	defer mu.Unlock()
	if loopCtx != nil {
		cancelLoopCtx()
	}
}

/*
Pause twenty-twenty-twenty.

This will pause the current twenty-twenty-twenty execution by [Settings]'s
'Pause' duration using [pkg/time.NewTimer].

The callback function in 'timerCallbackPre' parameter will be called once the
timer finishes.

The callback function in 'timerCallbackPos' parameter will be called once
twenty-twenty-twenty is resumed.
*/
func Pause(
	ctx context.Context,
	settings *Settings,
	optional Optional,
	timerCallbackPre func(),
	timerCallbackPos func(),
) {
	log.Printf("Pausing twenty-twenty-twenty for %.2f hour(s)...\n", settings.Pause.Hours())
	Stop() // cancelling current twenty-twenty-twenty goroutine
	timer := time.NewTimer(settings.Pause)

	select {
	case <-timer.C:
		log.Println("Resuming twenty-twenty-twenty...")
		if timerCallbackPre != nil {
			timerCallbackPre()
		}
		// need to start a new instance before calling the blocking
		// SendWithDuration(), otherwise if the user call Pause() again,
		// we are going to call Stop() in the previous loop
		Start(settings, optional)
		err := notification.SendWithDuration(
			ctx,
			&settings.Duration,
			&settings.Sound,
			"Resuming 20-20-20",
			fmt.Sprintf("You will see a notification every %.f minutes(s)", settings.Frequency.Minutes()),
		)
		if err != nil {
			log.Fatalf("Error while resuming notification: %v. Exiting...\n", err)
		}
		if timerCallbackPos != nil {
			timerCallbackPos()
		}
	case <-ctx.Done():
		log.Println("Cancelling twenty-twenty-twenty pause...")
	}
}

func loop(ctx context.Context, settings *Settings) {
	ticker := time.NewTicker(settings.Frequency)
	for {
		select {
		case <-ticker.C:
			log.Println("Sending notification...")
			go sound.SuspendAfter(settings.Duration * 2)
			err := notification.SendWithDuration(
				loopCtx,
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
			return
		}
	}
}
