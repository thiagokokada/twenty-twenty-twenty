package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"gioui.org/app"
	"gioui.org/x/notify"
)

type flags struct {
	durationInSec  uint
	frequencyInMin uint64
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
	flag.Parse()

	return flags{durationInSec: *durationInSec, frequencyInMin: *frequencyInMin}
}

func cancelNotificationAfter(notification notify.Notification, after time.Duration) {
	time.Sleep(after)
	err := notification.Cancel()
	if err != nil {
		fmt.Printf("Error while cancelling notification: %v\n", err)
	}
}

func main() {
	flags := parseFlags()

	notifier, err := notify.NewNotifier()
	if err != nil {
		log.Fatalf("Error while creating a notifier: %v\n", err)
	}

	notification, err := notifier.CreateNotification(
		"Starting 20-20-20",
		fmt.Sprintf("You will see a notification every %d minutes(s)", flags.frequencyInMin),
	)
	if err != nil {
		log.Fatalf("Error while sending test notification: %v\n", err)
	}
	go cancelNotificationAfter(notification, time.Duration(flags.durationInSec)*time.Second)

	ticker := time.NewTicker(time.Duration(flags.frequencyInMin) * time.Minute)
	fmt.Printf("Running twenty-twenty-twenty every %d minute(s)...\n", flags.frequencyInMin)
	go func() {
		for {
			<-ticker.C
			log.Println("Sending notification...")
			go func() {
				notification, err := notifier.CreateNotification(
					"Time to rest your eyes",
					fmt.Sprintf("Look at 20 feet (~6 meters) away for %d seconds", flags.durationInSec),
				)
				if err != nil {
					log.Printf("Error while sending notification: %v\n", err)
					return
				}
				go cancelNotificationAfter(notification, time.Duration(flags.durationInSec)*time.Second)
			}()
		}
	}()

	app.Main()
}
