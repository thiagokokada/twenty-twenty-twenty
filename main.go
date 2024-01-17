package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gioui.org/x/notify"
)

var (
	version           = "development"
	ctx               context.Context
	ctxCancel         context.CancelFunc
	duration          = new(time.Duration)
	frequency         = new(time.Duration)
	notificationSound = new(bool)
	notifier          notify.Notifier
)

type flags struct {
	disableSound   bool
	durationInSec  uint
	frequencyInMin float64
	version        bool
}

func parseFlags() flags {
	durationInSec := flag.Uint(
		"duration",
		20,
		"how long each pause should be in seconds",
	)
	frequencyInMin := flag.Float64(
		"frequency",
		20,
		"how often the pause should be in minutes",
	)
	version := flag.Bool(
		"version",
		false,
		"print program version and exit",
	)
	disableSound := new(bool)
	if notificationSoundEnabled {
		disableSound = flag.Bool(
			"disable-sound",
			false,
			"disable notification sound",
		)
	}
	flag.Parse()

	return flags{
		disableSound:   *disableSound,
		durationInSec:  *durationInSec,
		frequencyInMin: *frequencyInMin,
		version:        *version,
	}
}

func sendNotification(
	notifier notify.Notifier,
	title string,
	text string,
	notificationSound *bool,
) notify.Notification {
	if *notificationSound {
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
	notification notify.Notification,
	after *time.Duration,
	notificationSound *bool,
) {
	if notification == nil {
		return
	}
	time.Sleep(*after)

	if *notificationSound {
		playCancelNotificationSound()
	}

	err := notification.Cancel()
	if err != nil {
		fmt.Printf("Error while cancelling notification: %v\n", err)
	}
}

func twentyTwentyTwenty(
	ctx context.Context,
	notifier notify.Notifier,
	duration *time.Duration,
	frequency *time.Duration,
	notificationSound *bool,
) {
	ticker := time.NewTicker(*frequency)
	for {
		select {
		case <-ticker.C:
			go func() {
				log.Println("Sending notification...")
				notification := sendNotification(
					notifier,
					"Time to rest your eyes",
					fmt.Sprintf("Look at 20 feet (~6 meters) away for %.f seconds", duration.Seconds()),
					notificationSound,
				)
				go cancelNotificationAfter(notification, duration, notificationSound)
			}()
		case <-ctx.Done():
			log.Println("Disabling twenty-twenty-twenty...")
			return
		}
	}
}

func runTwentyTwentyTwenty() {
	if notificationSoundEnabled {
		log.Printf(
			"Running twenty-twenty-twenty every %.1f minute(s), with %.f second(s) duration and sound set to %t...\n",
			frequency.Minutes(),
			duration.Seconds(),
			*notificationSound,
		)
	} else {
		log.Printf(
			"Running twenty-twenty-twenty every %.1f minute(s), with %.f second(s) duration...\n",
			frequency.Minutes(),
			duration.Seconds(),
		)
	}

	ctx, ctxCancel = context.WithCancel(context.Background())
	go twentyTwentyTwenty(ctx, notifier, duration, frequency, notificationSound)
}

func main() {
	flags := parseFlags()
	if flags.version {
		fmt.Println(version)
		os.Exit(0)
	}

	*duration = time.Duration(flags.durationInSec) * time.Second
	*frequency = time.Duration(flags.frequencyInMin * float64(time.Minute))
	*notificationSound = notificationSoundEnabled && !flags.disableSound
	var err error

	// only init Beep if notification sound is enabled, otherwise we will cause
	// unnecessary noise in the speakers (and also increased memory usage)
	if *notificationSound {
		err = initNotification()
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
		fmt.Sprintf("You will see a notification every %.f minutes(s)", frequency.Minutes()),
		notificationSound,
	)
	if notification == nil {
		log.Fatalf("Test notification failed, exiting...")
	}
	go cancelNotificationAfter(notification, duration, notificationSound)

	runTwentyTwentyTwenty()
	loop()
}
