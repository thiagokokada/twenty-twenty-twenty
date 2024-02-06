//go:build !nosystray
// +build !nosystray

package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"fyne.io/systray"
)

const systrayEnabled bool = true

func withMutex(mu *sync.Mutex, f func()) {
	mu.Lock()
	defer mu.Unlock()
	f()
}

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
	if twenty.Optional.Sound {
		mSound = systray.AddMenuItemCheckbox("Sound", "Enable notification sound", twenty.Settings.Sound)
	}
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	var pauseCtx context.Context
	var cancelPauseCtx context.CancelFunc
	var mu sync.Mutex

	for {
		select {
		case <-mEnabled.ClickedCh:
			if mEnabled.Checked() {
				slog.DebugContext(twenty.Ctx(), "Enable button unchecked")
				twenty.Stop()

				withMutex(&mu, func() {
					mEnabled.Uncheck()
					mPause.Disable()
				})
			} else {
				slog.DebugContext(twenty.Ctx(), "Enable button checked")
				twenty.Start()

				withMutex(&mu, func() {
					mEnabled.Check()
					mPause.Enable()
				})
			}
		case <-mPause.ClickedCh:
			if pauseCtx != nil {
				slog.DebugContext(twenty.Ctx(), "Cancelling current pause")
				cancelPauseCtx()
			}
			if mPause.Checked() {
				slog.DebugContext(twenty.Ctx(), "Pause button unchecked")
				twenty.Start()

				withMutex(&mu, func() {
					mEnabled.Enable()
					mPause.Uncheck()
				})
			} else {
				slog.DebugContext(twenty.Ctx(), "Pause button checked")
				pauseCtx, cancelPauseCtx = context.WithCancel(context.Background())
				go func() {
					defer cancelPauseCtx()
					twenty.Pause(
						pauseCtx,
						func() {
							withMutex(&mu, func() {
								slog.DebugContext(twenty.Ctx(), "Calling pause callback")
								mEnabled.Enable()
								mPause.Uncheck()
							})
						},
						nil,
					)
				}()

				withMutex(&mu, func() {
					mEnabled.Disable()
					mPause.Check()
				})
			}
		case <-mSound.ClickedCh:
			if mSound.Checked() {
				slog.DebugContext(twenty.Ctx(), "Sound button unchecked")
				twenty.Settings.Sound = false

				withMutex(&mu, func() { mSound.Uncheck() })
			} else {
				slog.DebugContext(twenty.Ctx(), "Sound button checked")
				twenty.Settings.Sound = true

				withMutex(&mu, func() { mSound.Check() })
			}
		case <-mQuit.ClickedCh:
			slog.DebugContext(twenty.Ctx(), "Quit button clicked")
			systray.Quit()
		}
	}
}

func onExit() {
	twenty.Stop()
}
