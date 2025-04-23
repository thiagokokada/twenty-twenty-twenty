package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/thiagokokada/twenty-twenty-twenty/internal/core"
	"github.com/thiagokokada/twenty-twenty-twenty/internal/ctxlog"
	"github.com/thiagokokada/twenty-twenty-twenty/internal/notification"
	"github.com/thiagokokada/twenty-twenty-twenty/internal/sound"
)

var (
	version = "development"
	twenty  *core.TwentyTwentyTwenty
)

func main() {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)
	handler := &ctxlog.ContextHandler{
		Handler: logHandler{
			Handler:        slog.Default().Handler(),
			HandlerOptions: &slog.HandlerOptions{Level: lvl},
		},
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)
	// https://github.com/golang/go/issues/61892#issuecomment-1675123776
	// https://groups.google.com/g/golang-nuts/c/aJPXT2NF-Lc/m/rU0QayKwAQAJ
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime)

	features := core.Features{Sound: sound.Enabled, Systray: systrayEnabled}
	settings := core.ParseFlags(os.Args[0], os.Args[1:], version, features)

	if settings.Verbose {
		lvl.Set(slog.LevelDebug)
	}

	if features.Sound {
		err := sound.Init(!settings.Sound)
		if err != nil {
			log.Printf("Error while initialising sound: %v\n", err)
			log.Println("Disabling sound")
			features.Sound = false
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
	twenty = core.New(features, settings)
	// we need to start notification cancellation in a goroutine to show the
	// systray as soon as possible (since it depends on the loop() call), but we
	// also need to give it access to the core.Ctx to cancel it if necessary
	twenty.Start()
	go func() {
		if features.Sound {
			// wait the 1.5x of duration so we have some time for the sounds to
			// finish playing
			go sound.SuspendAfter(min(settings.Duration*3/2, settings.Frequency))
		}
		err := notification.CancelAfter(
			twenty.Ctx(),
			sentNotification,
			&twenty.Settings.Duration,
			&twenty.Settings.Sound,
		)
		if err != nil {
			log.Printf("Test notification cancel failed: %v\n", err)
		}
	}()

	loop()
}
