//go:build !nosystray
// +build !nosystray

package main

import (
	"context"
	"fmt"
	"log"

	"fyne.io/systray"

	"github.com/thiagokokada/twenty-twenty-twenty/core"
	snd "github.com/thiagokokada/twenty-twenty-twenty/sound"
)

const systrayEnabled bool = true

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
				go core.Pause(ctx, notifier, &settings,
					func() { mEnabled.Enable(); mPause.Uncheck(); ctxCancel() },
					func() { ctxCancel() },
				)

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
