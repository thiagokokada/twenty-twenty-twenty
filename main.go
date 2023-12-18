package main

import (
	_ "embed"
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

func main() {
	flags := parseFlags()

	notifier, err := notify.NewNotifier()
	if err != nil {
		log.Fatalf("Error while creating a notifier: %v\n", err)
	}

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
				time.Sleep(time.Duration(flags.durationInSec) * time.Second)
				notification.Cancel()
			}()
		}
	}()

	app.Main()
}
