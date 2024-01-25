package main

import (
	"fmt"
	"log"
	"os"

	"github.com/thiagokokada/twenty-twenty-twenty/internal/core"
	"github.com/thiagokokada/twenty-twenty-twenty/internal/notification"
	"github.com/thiagokokada/twenty-twenty-twenty/internal/sound"
)

var (
	version  = "development"
	optional core.Optional
	settings core.Settings
)

func main() {
	optional = core.Optional{Sound: sound.Enabled, Systray: systrayEnabled}
	settings = core.ParseFlags(os.Args[0], os.Args[1:], version, optional)

	if optional.Sound {
		err := sound.Init(!settings.Sound)
		if err != nil {
			log.Printf("Error while initialising sound: %v.\n", err)
			log.Println("Disabling sound...")
			optional.Sound = false
			settings.Sound = false
		}
	}

	sentNotification, err := notification.Send(
		&settings.Sound,
		"Starting 20-20-20",
		fmt.Sprintf("You will see a notification every %.f minutes(s)", settings.Frequency.Minutes()),
	)
	if err != nil {
		log.Fatalf("Test notification failed: %v. Exiting...", err)
	}
	// we need to start notification cancellation in a goroutine to show the
	// systray as soon as possible (since it depends on the loop() call), but we
	// also need to give it access to the core.Ctx to cancel it if necessary
	core.Start(&settings, optional)
	go func() {
		err := notification.CancelAfter(core.Ctx(), sentNotification, &settings.Duration, &settings.Sound)
		if err != nil {
			log.Printf("Test notification cancel failed: %v\n", err)
		}

		// wait the double of duration so we can make sure every sound finished
		// playing
		err = sound.SuspendAfter(settings.Duration * 2)
		if err != nil {
			log.Printf("Error while suspending speaker: %v\n", err)
		}
	}()

	loop()
}
