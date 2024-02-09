package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/jba/slog/handlers/loghandler"
	"github.com/thiagokokada/twenty-twenty-twenty/internal/core"
	"github.com/thiagokokada/twenty-twenty-twenty/internal/ctxlog"
	"github.com/thiagokokada/twenty-twenty-twenty/internal/notification"
	"github.com/thiagokokada/twenty-twenty-twenty/internal/sound"
)

var (
	ctx       context.Context
	ctxCancel context.CancelFunc
	twenty    *core.TwentyTwentyTwenty
	version   = "development"
)

func main() {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)
	handler := &ctxlog.ContextHandler{
		Handler: loghandler.New(
			os.Stdout,
			&slog.HandlerOptions{Level: lvl},
		),
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)

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
	ctx, ctxCancel = context.WithCancel(context.Background())
	go func() {
		twenty.Start(ctx)
		if features.Sound {
			// wait the 1.5x of duration so we have some time for the sounds to
			// finish playing
			go sound.SuspendAfter(min(settings.Duration*3/2, settings.Frequency))
		}
		err := notification.CancelAfter(
			ctx,
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
