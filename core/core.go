package core

import (
	"context"
	"fmt"
	"log"
	"time"

	"gioui.org/x/notify"
	n "github.com/thiagokokada/twenty-twenty-twenty/notification"
	s "github.com/thiagokokada/twenty-twenty-twenty/settings"
	snd "github.com/thiagokokada/twenty-twenty-twenty/sound"
)

var (
	Ctx    context.Context
	Cancel context.CancelFunc
)

func twentyTwentyTwenty(
	ctx context.Context,
	notifier notify.Notifier,
	settings *s.Settings,
) {
	ticker := time.NewTicker(settings.Frequency)
	cancelCtx, cancelCtxCancel := context.WithCancel(context.Background())
	for {
		select {
		case <-ticker.C:
			go func() {
				log.Println("Sending notification...")
				notification := n.Send(
					notifier,
					"Time to rest your eyes",
					fmt.Sprintf("Look at 20 feet (~6 meters) away for %.f seconds", settings.Duration.Seconds()),
					&settings.Sound,
				)
				go n.CancelAfter(cancelCtx, notification, &settings.Duration, &settings.Sound)
			}()
		case <-ctx.Done():
			log.Println("Disabling twenty-twenty-twenty...")
			cancelCtxCancel()
			return
		}
	}
}

func Start(
	notifier notify.Notifier,
	settings *s.Settings,
) {
	if snd.Enabled {
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

	Ctx, Cancel = context.WithCancel(context.Background())
	go twentyTwentyTwenty(Ctx, notifier, settings)
}
