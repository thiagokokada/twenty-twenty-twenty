package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gioui.org/x/notify"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/vorbis"
)

var (
	//go:embed notification.ogg
	NotificationSound embed.FS
	Buffer            *beep.Buffer
	Version           = "development"
)

type flags struct {
	durationInSec  uint
	frequencyInMin uint64
	version        bool
}

func parseFlags() flags {
	durationInSec := flag.Uint(
		"duration",
		20,
		"how long each pause should be in seconds",
	)
	frequencyInMin := flag.Uint64(
		"frequency",
		20,
		"how often the pause should be in minutes",
	)
	version := flag.Bool(
		"version",
		false,
		"print program version and exit",
	)
	flag.Parse()

	return flags{
		durationInSec:  *durationInSec,
		frequencyInMin: *frequencyInMin,
		version:        *version,
	}
}

func playNotificationSound() {
	done := make(chan bool)
	speaker.Play(
		beep.Seq(Buffer.Streamer(0, Buffer.Len())),
		beep.Callback(func() { done <- true }),
	)
	<-done
}

func sendNotification(notifier notify.Notifier, title string, text string) notify.Notification {
	notification, err := notifier.CreateNotification(title, text)
	if err != nil {
		log.Printf("Error while sending notification: %v\n", err)
		return nil
	}
	playNotificationSound()
	return notification
}

func cancelNotificationAfter(notification notify.Notification, after time.Duration) {
	if notification == nil {
		return
	}

	time.Sleep(after)
	err := notification.Cancel()
	if err != nil {
		fmt.Printf("Error while cancelling notification: %v\n", err)
	}
}

func twentyTwentyTwenty(notifier notify.Notifier, duration time.Duration, frequency time.Duration) {
	ticker := time.NewTicker(frequency)
	for {
		<-ticker.C
		go func() {
			log.Println("Sending notification...")
			notification := sendNotification(
				notifier,
				"Time to rest your eyes",
				fmt.Sprintf("Look at 20 feet (~6 meters) away for %.f seconds", duration.Seconds()),
			)
			go cancelNotificationAfter(notification, duration)
		}()
	}
}

func init() {
	f, err := NotificationSound.Open("notification.ogg")
	if err != nil {
		log.Fatalf("Failed to load notification sound: %v\n", err)
	}

	streamer, format, err := vorbis.Decode(f)
	if err != nil {
		log.Fatalf("Failed to decode the notification sound: %v\n", err)
	}
	Buffer = beep.NewBuffer(format)
	Buffer.Append(streamer)

	// 1s/4 = 250ms of lag, good enough for this use case
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/4))
}

func main() {
	flags := parseFlags()
	if flags.version {
		fmt.Println(Version)
		os.Exit(0)
	}

	duration := time.Duration(flags.durationInSec) * time.Second
	frequency := time.Duration(flags.frequencyInMin) * time.Minute

	notifier, err := notify.NewNotifier()
	if err != nil {
		log.Fatalf("Error while creating a notifier: %v\n", err)
	}

	notification := sendNotification(
		notifier,
		"Starting 20-20-20",
		fmt.Sprintf("You will see a notification every %.f minutes(s)", frequency.Minutes()),
	)
	if notification == nil {
		log.Fatalf("Test notification failed, exiting...")
	}
	go cancelNotificationAfter(notification, duration)

	fmt.Printf("Running twenty-twenty-twenty every %.f minute(s)...\n", frequency.Minutes())
	go twentyTwentyTwenty(notifier, duration, frequency)
	loop()
}
