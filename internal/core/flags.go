package core

import (
	"flag"
	"fmt"
	"os"
	"time"
)

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
