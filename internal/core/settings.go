package core

import (
	"flag"
	"fmt"
	"os"
	"time"
)

/*
Settings struct.

'Duration' will be the duration of each notification. For example, if is 20
seconds, it means that each notification will stay by 20 seconds.

'Frequency' is how often each notification will be shown. For example, if it is
20 minutes, a new notification will appear at every 20 minutes.

'Pause' is the duration of the pause. For example, if it is 1 hour, we will
disable notifications for 1 hour.

'Sound' enables or disables sound every time a notification is shown.
*/
type Settings struct {
	Duration  time.Duration
	Frequency time.Duration
	Pause     time.Duration
	Sound     bool
	Verbose   bool
}

func ParseFlags(
	progname string,
	args []string,
	version string,
	optional Optional,
) Settings {
	flags := flag.NewFlagSet(progname, flag.ExitOnError)
	durationInSec := flags.Uint(
		"duration",
		20,
		"how long each pause should be in seconds",
	)
	frequencyInSec := flags.Uint(
		"frequency",
		20*60,
		"how often the pause should be in seconds",
	)
	pauseInSec := new(uint)
	if optional.Systray {
		pauseInSec = flags.Uint(
			"pause",
			60*60,
			"how long the pause (from systray) should be in seconds",
		)
	}
	disableSound := new(bool)
	if optional.Sound {
		disableSound = flags.Bool(
			"disable-sound",
			false,
			"disable notification sound",
		)
	}
	showVersion := flags.Bool(
		"version",
		false,
		"print program version and exit",
	)
	verbose := flags.Bool(
		"verbose",
		false,
		"enable verbose logging",
	)
	// using flag.ExitOnError, so no error will ever be returned
	_ = flags.Parse(args)

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	return Settings{
		Duration:  time.Duration(*durationInSec) * time.Second,
		Frequency: time.Duration(*frequencyInSec) * time.Second,
		Pause:     time.Duration(*pauseInSec) * time.Second,
		Sound:     optional.Sound && !*disableSound,
		Verbose:   *verbose,
	}
}
