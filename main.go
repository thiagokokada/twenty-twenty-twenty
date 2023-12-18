package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gen2brain/beeep"
)

type flags struct {
	durationInSec  uint
	frequencyInMin uint64
}

//go:embed eye-solid.svg
var iconBytes []byte

func loadIcon() *os.File {
	icon, err := os.CreateTemp("", "twenty-twenty-twenty-icon-")
	if err != nil {
		log.Fatalf("Error while creating temporary file: %v", err)
	}
	icon.Write(iconBytes)
	return icon
}

func parseFlags() flags {
	durationInSec := flag.Uint(
		"duration",
		20,
		"how long to show the notification in seconds (does not work in macOS)",
	)
	frequencyInMin := flag.Uint64(
		"frequency",
		20,
		"how often to show the notification in minutes",
	)
	flag.Parse()

	return flags{durationInSec: *durationInSec, frequencyInMin: *frequencyInMin}
}

func initBeeep(durationInSec uint) {
	const MS_IN_SEC = 1000

	err := beeep.Beep(beeep.DefaultFreq, int(durationInSec)*MS_IN_SEC)
	if err != nil {
		log.Fatalf("Error during beeep init: %v\n", err)
	}
}

func main() {
	icon := loadIcon()
	defer icon.Close()
	defer os.Remove(icon.Name())

	flags := parseFlags()
	initBeeep(flags.durationInSec)

	ticker := time.NewTicker(time.Duration(flags.frequencyInMin) * time.Minute)

	fmt.Printf("Running twenty-twenty-twenty every %d minute(s)...\n", flags.frequencyInMin)

	for {
		<-ticker.C
		log.Println("Sending notification...")
		err := beeep.Alert(
			"Time to rest your eyes",
			fmt.Sprintf("Look at 20 feet (~6 meters) away for %d seconds", flags.durationInSec),
			icon.Name(),
		)
		if err != nil {
			log.Fatalf("Error during beeep alert: %v\n", err)
		}
	}
}
