//go:build !nosystray
// +build !nosystray

package main

import (
	"context"
	"fmt"
	"log"

	"fyne.io/systray"

	"github.com/thiagokokada/twenty-twenty-twenty/core"
	"github.com/thiagokokada/twenty-twenty-twenty/sound"
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
	if sound.Enabled {
		mSound = systray.AddMenuItemCheckbox("Sound", "Enable notification sound", settings.Sound)
	}
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	var pauseCtx context.Context
	var cancelPauseCtx context.CancelFunc

	for {
		select {
		case <-mEnabled.ClickedCh:
			if mEnabled.Checked() {
				core.Stop()

				mEnabled.Uncheck()
				mPause.Disable()
			} else {
				core.Start(&settings, optional)

				mEnabled.Check()
				mPause.Enable()
			}
		case <-mPause.ClickedCh:
			if mPause.Checked() {
				cancelPauseCtx() // cancel the current pause if it is running
				core.Start(&settings, optional)

				mEnabled.Enable()
				mPause.Uncheck()
			} else {
				pauseCtx, cancelPauseCtx = context.WithCancel(context.Background())
				go func() {
					defer cancelPauseCtx()
					core.Pause(
						pauseCtx, &settings, optional,
						func() { mPause.Disable() }, // blocking pause button to avoid concurrency issue
						func() { mEnabled.Enable(); mPause.Uncheck(); mPause.Enable() },
					)
				}()

				mEnabled.Disable()
				mPause.Check()
			}
		case <-mSound.ClickedCh:
			if mSound.Checked() {
				settings.Sound = false

				mSound.Uncheck()
			} else {
				err := sound.Init()
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
