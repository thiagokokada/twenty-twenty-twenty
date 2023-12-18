package main

import (
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
