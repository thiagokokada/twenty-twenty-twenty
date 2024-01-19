//go:build !nosystray
// +build !nosystray

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"fyne.io/systray"

	"github.com/thiagokokada/twenty-twenty-twenty/core"
	ntf "github.com/thiagokokada/twenty-twenty-twenty/notification"
	snd "github.com/thiagokokada/twenty-twenty-twenty/sound"
)

const systrayEnabled bool = true

type menuItems struct {
	mEnabled *systray.MenuItem
	mPause   *systray.MenuItem
	mQuit    *systray.MenuItem
	mSound   *systray.MenuItem
}

func resumeTwentyTwentyTwentyAfter(
	ctx context.Context,
	ctxCancel context.CancelFunc,
	settings *core.Settings,
	menu *menuItems,
) {
	log.Printf("Pausing twenty-twenty-twenty for %.f hour...\n", settings.Pause.Hours())
	core.Cancel() // cancelling current twenty-twenty-twenty goroutine
	timer := time.NewTimer(settings.Pause)
	cancelCtx, cancelCtxCancel := context.WithCancel(context.Background())

	select {
	case <-timer.C:
		notification := ntf.Send(
			notifier,
			"Resuming 20-20-20",
			fmt.Sprintf("You will see a notification every %.f minutes(s)", settings.Frequency.Minutes()),
			&settings.Sound,
		)
		if notification == nil {
			log.Printf("Resume notification failed...")
		}
		go ntf.CancelAfter(cancelCtx, notification, &settings.Duration, &settings.Sound)
		core.Start(notifier, settings)

		menu.mEnabled.Enable()
		menu.mPause.Uncheck()
	case <-ctx.Done():
	}
	cancelCtxCancel()
	ctxCancel() // make sure the current context is closed
}

func onReady() {
	systray.SetIcon(systrayIcon)
	systray.SetTooltip("TwentyTwentyTwenty")
	mEnabled := systray.AddMenuItemCheckbox("Enabled", "Enable twenty-twenty-twenty", true)
	mPause := systray.AddMenuItemCheckbox(
		fmt.Sprintf("Pause for %.f hour", settings.Pause.Hours()),
		fmt.Sprintf("Pause twenty-twenty-twenty for %.f hour", settings.Pause.Hours()),
		false,
	)
	mSound := new(systray.MenuItem)
	if snd.Enabled {
		mSound = systray.AddMenuItemCheckbox("Sound", "Enable notification sound", settings.Sound)
	}
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	menu := menuItems{mEnabled, mPause, mSound, mQuit}

	var ctx context.Context
	var ctxCancel context.CancelFunc

	for {
		select {
		case <-mEnabled.ClickedCh:
			if mEnabled.Checked() {
				core.Cancel()

				mEnabled.Uncheck()
				mPause.Disable()
			} else {
				core.Start(notifier, &settings)

				mEnabled.Check()
				mPause.Enable()
			}
		case <-mPause.ClickedCh:
			if mPause.Checked() {
				core.Cancel() // make sure the current twenty-twenty-twenty goroutine stopped
				ctxCancel()   // cancel the current pause if it is running
				core.Start(notifier, &settings)

				mEnabled.Enable()
				mPause.Uncheck()
			} else {
				ctx, ctxCancel = context.WithCancel(context.Background())
				go resumeTwentyTwentyTwentyAfter(ctx, ctxCancel, &settings, &menu)

				mEnabled.Disable()
				mPause.Check()
			}
		case <-mSound.ClickedCh:
			if mSound.Checked() {
				settings.Sound = false

				mSound.Uncheck()
			} else {
				err := snd.Init()
				if err != nil {
					log.Fatalf("Error while initialising sound: %v\n", err)
				}
				settings.Sound = true

				mSound.Check()
			}
		case <-mQuit.ClickedCh:
			systray.Quit()
		}
	}
}

func onExit() {}
