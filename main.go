package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/x/notify"

	"github.com/thiagokokada/twenty-twenty-twenty/core"
	ntf "github.com/thiagokokada/twenty-twenty-twenty/notification"
	"github.com/thiagokokada/twenty-twenty-twenty/sound"
)

var (
	version  = "development"
	notifier notify.Notifier
	optional core.Optional
	settings core.Settings
)

func main() {
	optional = core.Optional{Sound: sound.Enabled, Systray: systrayEnabled}
	settings = core.ParseFlags(os.Args[0], os.Args[1:], version, optional)
	var err error

	// only init Beep if notification sound is enabled, otherwise we will cause
	// unnecessary noise in the speakers (and also increased memory usage)
	if settings.Sound {
		err = sound.Init()
		if err != nil {
			log.Fatalf("Error while initialising sound: %v\n", err)
		}
	}

	notifier, err = notify.NewNotifier()
	if err != nil {
		log.Fatalf("Error while creating a notifier: %v\n", err)
	}

	notification, err := ntf.Send(
		notifier,
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
	core.Start(notifier, &settings, optional)
	go ntf.CancelAfter(core.Ctx, notification, &settings.Duration, &settings.Sound)

	loop()
}
