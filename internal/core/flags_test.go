package core

import (
	"strings"
	"testing"
	"time"

	"github.com/thiagokokada/twenty-twenty-twenty/internal/assert"
)

func TestParseFlags(t *testing.T) {
	const progname = "twenty-twenty-twenty"
	const version = "test"

	// always return false for sound if disabled
	settings := ParseFlags(progname, []string{}, version, Features{Sound: false, Systray: false})
	assert.Equal(t, settings.Sound, false)

	var tests = []struct {
		args     []string
		settings Settings
	}{
		{[]string{},
			Settings{Duration: time.Second * 20, Frequency: time.Minute * 20, Pause: time.Hour, Sound: true, Verbose: false}},
		{[]string{"-duration", "10", "-frequency", "600", "-pause", "1800", "-disable-sound", "-verbose"},
			Settings{Duration: time.Second * 10, Frequency: time.Minute * 10, Pause: time.Minute * 30, Sound: false, Verbose: true}},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			settings := ParseFlags(progname, tt.args, version, Features{Sound: true, Systray: true})
			assert.Equal(t, settings, tt.settings)
		})
	}
}
