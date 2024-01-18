package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"gioui.org/x/notify"

	s "github.com/thiagokokada/twenty-twenty-twenty/settings"
)

var (
	version       = "development"
	mainCtx       context.Context
	mainCtxCancel context.CancelFunc
	notifier      notify.Notifier
	settings      s.Settings
)

func sendNotification(
	notifier notify.Notifier,
	title string,
	text string,
	sound *bool,
) notify.Notification {
	if *sound {
		playSendNotificationSound()
	}

	notification, err := notifier.CreateNotification(title, text)
	if err != nil {
		log.Printf("Error while sending notification: %v\n", err)
		return nil
	}
	return notification
}

func cancelNotificationAfter(
	ctx context.Context,
	after *time.Duration,
	notification notify.Notification,
) {
	if notification == nil {
		return
	}

	timer := time.NewTimer(*after)
	select {
	case <-timer.C:
		if settings.Sound {
			playCancelNotificationSound()
		}
	case <-ctx.Done(): // avoid playing notification sound if we cancel the context
	}
	err := notification.Cancel()
	if err != nil {
		log.Printf("Error while cancelling notification: %v\n", err)
	}
}

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
				notification := sendNotification(
					notifier,
					"Time to rest your eyes",
					fmt.Sprintf("Look at 20 feet (~6 meters) away for %.f seconds", settings.Duration.Seconds()),
					&settings.Sound,
				)
				go cancelNotificationAfter(cancelCtx, &settings.Duration, notification)
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
	if notificationSoundEnabled {
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
	settings = s.ParseFlags(version, systrayEnabled, notificationSoundEnabled)
	var err error

	// only init Beep if notification sound is enabled, otherwise we will cause
	// unnecessary noise in the speakers (and also increased memory usage)
	if settings.Sound {
		err = initSound()
		if err != nil {
			log.Fatalf("Error while initialising sound: %v\n", err)
		}
	}

	notifier, err = notify.NewNotifier()
	if err != nil {
		log.Fatalf("Error while creating a notifier: %v\n", err)
	}

	notification := sendNotification(
		notifier,
		"Starting 20-20-20",
		fmt.Sprintf("You will see a notification every %.f minutes(s)", settings.Frequency.Minutes()),
		&settings.Sound,
	)
	if notification == nil {
		log.Fatalf("Test notification failed, exiting...")
	}
	go cancelNotificationAfter(context.Background(), &settings.Duration, notification)

	runTwentyTwentyTwenty(notifier, &settings)
	loop()
}
