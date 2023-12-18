package main

import (
	"math"
	"os"
	"testing"
)

func TestLoadIcon(t *testing.T) {
	icon := loadIcon()
	f, err := os.Open(icon.Name())
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		t.Error(err)
	}
	if stat.Size() <= 100 {
		t.Errorf("Icon should be at least 100 bytes, it was %d", stat.Size())
	}
}

func TestInitBeeep(t *testing.T) {
	initBeeep(1)
	initBeeep(10)
	initBeeep(100)
	initBeeep(1000)
	initBeeep(math.MaxUint64)
}
