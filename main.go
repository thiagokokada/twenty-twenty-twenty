package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"gioui.org/x/notify"

	s "github.com/thiagokokada/twenty-twenty-twenty/settings"
	snd "github.com/thiagokokada/twenty-twenty-twenty/sound"
	n "github.com/thiagokokada/twenty-twenty-twenty/notification"
)

var (
	version       = "development"
	mainCtx       context.Context
	mainCtxCancel context.CancelFunc
	notifier      notify.Notifier
	settings      s.Settings
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

func runTwentyTwentyTwenty(
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

	mainCtx, mainCtxCancel = context.WithCancel(context.Background())
	go twentyTwentyTwenty(mainCtx, notifier, settings)
}

func main() {
	settings = s.ParseFlags(version, systrayEnabled, snd.Enabled)
	var err error

	// only init Beep if notification sound is enabled, otherwise we will cause
	// unnecessary noise in the speakers (and also increased memory usage)
	if settings.Sound {
		err = snd.Init()
		if err != nil {
			log.Fatalf("Error while initialising sound: %v\n", err)
		}
	}

	notifier, err = notify.NewNotifier()
	if err != nil {
		log.Fatalf("Error while creating a notifier: %v\n", err)
	}

	notification := n.Send(
		notifier,
		"Starting 20-20-20",
		fmt.Sprintf("You will see a notification every %.f minutes(s)", settings.Frequency.Minutes()),
		&settings.Sound,
	)
	if notification == nil {
		log.Fatalf("Test notification failed, exiting...")
	}
	go n.CancelAfter(context.Background(), notification, &settings.Duration, &settings.Sound)

	runTwentyTwentyTwenty(notifier, &settings)
	loop()
}
