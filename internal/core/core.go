package core

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/thiagokokada/twenty-twenty-twenty/internal/ctxlog"
	"github.com/thiagokokada/twenty-twenty-twenty/internal/notification"
	"github.com/thiagokokada/twenty-twenty-twenty/internal/sound"
)

var loop atomic.Int64

/*
Create a new TwentyTwentyTwenty struct.
*/
func New(features Features, settings Settings) *TwentyTwentyTwenty {
	return &TwentyTwentyTwenty{Features: features, Settings: settings}
}

/*
Start twenty-twenty-twenty.

This will start the main twenty-twenty-twenty loop in a goroutine, so avoid
calling this function inside a goroutine.
*/
func (t *TwentyTwentyTwenty) Start(ctx context.Context) {
	if t.Features.Sound {
		log.Printf(
			"Running twenty-twenty-twenty every %.1f minute(s), with %.f second(s) duration and sound set to %t\n",
			t.Settings.Frequency.Minutes(),
			t.Settings.Duration.Seconds(),
			t.Settings.Sound,
		)
	} else {
		log.Printf(
			"Running twenty-twenty-twenty every %.1f minute(s), with %.f second(s) duration\n",
			t.Settings.Frequency.Minutes(),
			t.Settings.Duration.Seconds(),
		)
	}

	t.loop(ctxlog.AppendCtx(ctx, slog.Int64("loop", loop.Add(1))))
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
func (t *TwentyTwentyTwenty) Pause(
	ctx context.Context,
	timerCallbackPre func(),
	timerCallbackPos func(),
) {
	log.Printf("Pausing twenty-twenty-twenty for %.2f hour(s)\n", t.Settings.Pause.Hours())
	timer := time.NewTimer(t.Settings.Pause)
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
		go t.Start(ctx)
		err := notification.SendWithDuration(
			ctx,
			&t.Settings.Duration,
			&t.Settings.Sound,
			"Resuming 20-20-20",
			fmt.Sprintf("You will see a notification every %.f minutes(s)", t.Settings.Frequency.Minutes()),
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

func (t *TwentyTwentyTwenty) loop(ctx context.Context) {
	slog.DebugContext(ctx, "Starting new loop")
	ticker := time.NewTicker(t.Settings.Frequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Printf("Showing notification for %.f second(s)\n", t.Settings.Duration.Seconds())
			// wait 1.5x the duration so we have some time for the sounds to
			// finish playing
			if t.Features.Sound {
				go sound.SuspendAfter(min(t.Settings.Duration*3/2, t.Settings.Frequency))
			}
			err := notification.SendWithDuration(
				ctx,
				&t.Settings.Duration,
				&t.Settings.Sound,
				"Time to rest your eyes",
				fmt.Sprintf("Look at 20 feet (~6 meters) away for %.f seconds", t.Settings.Duration.Seconds()),
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
