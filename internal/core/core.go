package core

import (
	"context"
	"fmt"
	"log"
	"log/slog"
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
Optional struct.

This is used for features that are optional in the program, for example if sound
or systray are permanently disabled.
*/
type Optional struct {
	Sound   bool
	Systray bool
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
			"Running twenty-twenty-twenty every %.1f minute(s), with %.f second(s) duration and sound set to %t\n",
			settings.Frequency.Minutes(),
			settings.Duration.Seconds(),
			settings.Sound,
		)
	} else {
		log.Printf(
			"Running twenty-twenty-twenty every %.1f minute(s), with %.f second(s) duration\n",
			settings.Frequency.Minutes(),
			settings.Duration.Seconds(),
		)
	}
	Stop() // make sure we cancel the previous instance

	mu.Lock()
	defer mu.Unlock()
	loopCtx, cancelLoopCtx = context.WithCancel(context.Background())
	go loop(loopCtx, settings, optional)
}

/*
Stop twenty-twenty-twenty.
*/
func Stop() {
	mu.Lock()
	defer mu.Unlock()
	if loopCtx != nil {
		slog.Debug("Cancelling main loop context")
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
	log.Printf("Pausing twenty-twenty-twenty for %.2f hour(s)\n", settings.Pause.Hours())
	Stop() // cancelling current twenty-twenty-twenty goroutine
	timer := time.NewTimer(settings.Pause)
	defer timer.Stop()

	select {
	case <-timer.C:
		log.Println("Resuming twenty-twenty-twenty")
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
			log.Fatalf("Error while resuming notification: %v. Exiting\n", err)
		}
		if timerCallbackPos != nil {
			timerCallbackPos()
		}
	case <-ctx.Done():
		slog.DebugContext(ctx, "Cancelling twenty-twenty-twenty pause")
	}
}

func loop(ctx context.Context, settings *Settings, optional Optional) {
	slog.DebugContext(ctx, "Starting new loop")
	ticker := time.NewTicker(settings.Frequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Printf("Showing notification for %.f second(s)\n", settings.Duration.Seconds())
			// wait 1.5x the duration so we have some time for the sounds to
			// finish playing
			if optional.Sound {
				go sound.SuspendAfter(min(settings.Duration*3/2, settings.Frequency))
			}
			err := notification.SendWithDuration(
				ctx,
				&settings.Duration,
				&settings.Sound,
				"Time to rest your eyes",
				fmt.Sprintf("Look at 20 feet (~6 meters) away for %.f seconds", settings.Duration.Seconds()),
			)
			if err != nil {
				log.Printf("Error while sending notification: %v\n", err)
			}
		case <-ctx.Done():
			log.Println("Disabling twenty-twenty-twenty")
			return
		}
	}
}
