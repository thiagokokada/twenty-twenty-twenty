//go:build !nosystray
// +build !nosystray

package main

import (
	"context"
	"fmt"
	"log/slog"

	"fyne.io/systray"
)

const systrayEnabled bool = true

func onReady() {
	systray.SetIcon(systrayIcon)
	systray.SetTooltip("TwentyTwentyTwenty")
	mEnabled := systray.AddMenuItemCheckbox("Enabled", "Enable twenty-twenty-twenty", true)
	mPause := systray.AddMenuItemCheckbox(
		fmt.Sprintf("Pause for %.f hour", twenty.Settings.Pause.Hours()),
		fmt.Sprintf("Pause twenty-twenty-twenty for %.f hour", twenty.Settings.Pause.Hours()),
		false,
	)
	mSound := new(systray.MenuItem)
	if twenty.Features.Sound {
		mSound = systray.AddMenuItemCheckbox("Sound", "Enable notification sound", twenty.Settings.Sound)
	}
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// running `go run . -race` detects a bunch of data races here that are safe
	// to ignore since the access to the actual critical regions in systray are
	// protected with a mutex
	for {
		select {
		case <-mEnabled.ClickedCh:
			if mEnabled.Checked() {
				slog.DebugContext(ctx, "Enable button unchecked")
				ctxCancel()

				mEnabled.Uncheck()
				mPause.Disable()
			} else {
				slog.DebugContext(ctx, "Enable button checked")
				ctx, ctxCancel = context.WithCancel(context.Background())
				go twenty.Start(ctx)

				mEnabled.Check()
				mPause.Enable()
			}
		case <-mPause.ClickedCh:
			slog.DebugContext(ctx, "Cancelling current pause")
			ctxCancel()
			if mPause.Checked() {
				slog.DebugContext(ctx, "Pause button unchecked")
				ctx, ctxCancel = context.WithCancel(context.Background())
				go twenty.Start(ctx)

				mEnabled.Enable()
				mPause.Uncheck()
			} else {
				slog.DebugContext(ctx, "Pause button checked")
				ctx, ctxCancel = context.WithCancel(context.Background())
				go twenty.Pause(
					ctx,
					func() {
						slog.DebugContext(ctx, "Calling pause callback")
						mEnabled.Enable()
						mPause.Uncheck()
					},
					nil,
				)

				mEnabled.Disable()
				mPause.Check()
			}
		case <-mSound.ClickedCh:
			if mSound.Checked() {
				slog.DebugContext(ctx, "Sound button unchecked")
				twenty.Settings.Sound = false

				mSound.Uncheck()
			} else {
				slog.DebugContext(ctx, "Sound button checked")
				twenty.Settings.Sound = true

				mSound.Check()
			}
		case <-mQuit.ClickedCh:
			slog.DebugContext(ctx, "Quit button clicked")
			systray.Quit()
		}
	}
}

func onExit() {
	ctxCancel()
}
