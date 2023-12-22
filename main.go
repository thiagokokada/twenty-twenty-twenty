package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gioui.org/x/notify"
)

var Version = "development"

type flags struct {
	durationInSec  uint
	frequencyInMin uint64
	version        bool
}

func parseFlags() flags {
	durationInSec := flag.Uint(
		"duration",
		20,
		"how long each pause should be in seconds",
	)
	frequencyInMin := flag.Uint64(
		"frequency",
		20,
		"how often the pause should be in minutes",
	)
	version := flag.Bool(
		"version",
		false,
		"print program version and exit",
	)
	flag.Parse()

	return flags{
		durationInSec:  *durationInSec,
		frequencyInMin: *frequencyInMin,
		version:        *version,
	}
}

func sendNotification(notifier notify.Notifier, title string, text string) notify.Notification {
	notification, err := notifier.CreateNotification(title, text)
	if err != nil {
		log.Printf("Error while sending notification: %v\n", err)
		return nil
	}
	return notification
}

func cancelNotificationAfter(notification notify.Notification, after time.Duration) {
	if notification == nil {
		return
	}

	time.Sleep(after)
	err := notification.Cancel()
	if err != nil {
		fmt.Printf("Error while cancelling notification: %v\n", err)
	}
}

func twentyTwentyTwenty(notifier notify.Notifier, duration time.Duration, frequency time.Duration) {
	ticker := time.NewTicker(frequency)
	for {
		<-ticker.C
		go func() {
			log.Println("Sending notification...")
			notification := sendNotification(
				notifier,
				"Time to rest your eyes",
				fmt.Sprintf("Look at 20 feet (~6 meters) away for %.f seconds", duration.Seconds()),
			)
			go cancelNotificationAfter(notification, duration)
		}()
	}
}

func main() {
	flags := parseFlags()
	if flags.version {
		fmt.Println(Version)
		os.Exit(0)
	}

	duration := time.Duration(flags.durationInSec) * time.Second
	frequency := time.Duration(flags.frequencyInMin) * time.Minute

	notifier, err := notify.NewNotifier()
	if err != nil {
		log.Fatalf("Error while creating a notifier: %v\n", err)
	}

	notification := sendNotification(
		notifier,
		"Starting 20-20-20",
		fmt.Sprintf("You will see a notification every %.f minutes(s)", frequency.Minutes()),
	)
	if notification == nil {
		log.Fatalf("Test notification failed, exiting...")
	}
	go cancelNotificationAfter(notification, duration)

	fmt.Printf("Running twenty-twenty-twenty every %.f minute(s)...\n", frequency.Minutes())
	go twentyTwentyTwenty(notifier, duration, frequency)
	loop()
}
