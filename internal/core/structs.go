package core

import (
	"context"
	"sync"
	"time"
)

/*
TwentyTwentyTwenty struct.

Keeps the main state of the program.
*/
type TwentyTwentyTwenty struct {
	Optional Optional
	Settings Settings

	cancelLoopCtx context.CancelFunc
	loopCtx       context.Context
	mu            sync.Mutex
}

/*
Optional struct.

This is used for features that are optional in the program, for example if sound
or systray are permanently disabled.
*/
type Optional struct {
	Sound   bool
	Systray bool
}

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
