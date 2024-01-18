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
	version       = "development"
	mainCtx       context.Context
	mainCtxCancel context.CancelFunc
	notifier      notify.Notifier
	settings      appSettings
)

type appSettings struct {
	duration  time.Duration
	frequency time.Duration
	pause     time.Duration
	sound     bool
}

func parseFlags() appSettings {
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
	if notificationSoundEnabled {
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

	return appSettings{
		duration:  time.Duration(*durationInSec) * time.Second,
		frequency: time.Duration(*frequencyInSec) * time.Second,
		pause:     time.Duration(*pauseInSec) * time.Second,
		sound:     notificationSoundEnabled && !*disableSound,
	}
}

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
	notification notify.Notification,
	settings *appSettings,
) {
	if notification == nil {
		return
	}
	time.Sleep(settings.duration)

	if settings.sound {
		playCancelNotificationSound()
	}

	err := notification.Cancel()
	if err != nil {
		log.Printf("Error while cancelling notification: %v\n", err)
	}
}

func twentyTwentyTwenty(
	ctx context.Context,
	notifier notify.Notifier,
	settings *appSettings,
) {
	ticker := time.NewTicker(settings.frequency)
	for {
		select {
		case <-ticker.C:
			go func() {
				log.Println("Sending notification...")
				notification := sendNotification(
					notifier,
					"Time to rest your eyes",
					fmt.Sprintf("Look at 20 feet (~6 meters) away for %.f seconds", settings.duration.Seconds()),
					&settings.sound,
				)
				go cancelNotificationAfter(notification, settings)
			}()
		case <-ctx.Done():
			log.Println("Disabling twenty-twenty-twenty...")
			return
		}
	}
}

func runTwentyTwentyTwenty(
	notifier notify.Notifier,
	settings *appSettings,
) {
	if notificationSoundEnabled {
		log.Printf(
			"Running twenty-twenty-twenty every %.1f minute(s), with %.f second(s) duration and sound set to %t...\n",
			settings.frequency.Minutes(),
			settings.duration.Seconds(),
			settings.sound,
		)
	} else {
		log.Printf(
			"Running twenty-twenty-twenty every %.1f minute(s), with %.f second(s) duration...\n",
			settings.frequency.Minutes(),
			settings.duration.Seconds(),
		)
	}

	mainCtx, mainCtxCancel = context.WithCancel(context.Background())
	go twentyTwentyTwenty(mainCtx, notifier, settings)
}

func main() {
	settings = parseFlags()
	var err error

	// only init Beep if notification sound is enabled, otherwise we will cause
	// unnecessary noise in the speakers (and also increased memory usage)
	if settings.sound {
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
		fmt.Sprintf("You will see a notification every %.f minutes(s)", settings.frequency.Minutes()),
		&settings.sound,
	)
	if notification == nil {
		log.Fatalf("Test notification failed, exiting...")
	}
	go cancelNotificationAfter(notification, &settings)

	runTwentyTwentyTwenty(notifier, &settings)
	loop()
}
