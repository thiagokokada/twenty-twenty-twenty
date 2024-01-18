package settings

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type Settings struct {
	Duration  time.Duration
	Frequency time.Duration
	Pause     time.Duration
	Sound     bool
}

func ParseFlags(
	version string,
	systrayEnabled bool,
	notificationSoundEnabled bool,
) Settings {
	durationInSec := flag.Uint(
		"duration",
		20,
		"how long each pause should be in seconds",
	)
	frequencyInSec := flag.Uint(
		"frequency",
		20*60,
		"how often the pause should be in seconds",
	)
	pauseInSec := new(uint)
	if systrayEnabled {
		pauseInSec = flag.Uint(
			"pause",
			60*60,
			"how long the pause (from systray) should be in seconds",
		)
	}
	disableSound := new(bool)
	if notificationSoundEnabled {
		disableSound = flag.Bool(
			"disable-sound",
			false,
			"disable notification sound",
		)
	}
	showVersion := flag.Bool(
		"version",
		false,
		"print program version and exit",
	)
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	return Settings{
		Duration:  time.Duration(*durationInSec) * time.Second,
		Frequency: time.Duration(*frequencyInSec) * time.Second,
		Pause:     time.Duration(*pauseInSec) * time.Second,
		Sound:     notificationSoundEnabled && !*disableSound,
	}
}
