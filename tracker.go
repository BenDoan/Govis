package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"stathat.com/c/jconfig"
)

type Tracker struct {
	lastTime       time.Time
	lastWindow     string
	currentWindow  string
	interval       time.Duration
	minIdleTime    time.Duration
	ignorePatterns []interface{}
}

func (t *Tracker) Start(config *jconfig.Config) {
	minIdleTime, err := time.ParseDuration(strconv.Itoa(config.GetInt("MinIdleTime")) + "s")

	if err != nil {
		log.Error("Could not parse MinIdleTime: %v", err)
		t.minIdleTime = 5 * time.Second
	} else {
		t.minIdleTime = minIdleTime
	}

	t.ignorePatterns = config.GetArray("IgnorePatterns")

	t.interval = time.Duration(config.GetInt("TickInterval")) * time.Millisecond
	t.lastTime = time.Now()
	t.lastWindow = GetCurrentWindowName()

	t.StartTracking()
}

func (t *Tracker) StartTracking() {
	c := time.Tick(t.interval)
	for now := range c {
		t.currentWindow = GetCurrentWindowName()

		if t.currentWindow != t.lastWindow && GetIdleTime() < t.minIdleTime {
			t.PrintStatus()
		}

		t.lastWindow = t.currentWindow
		var _ = now
	}
}

func (t *Tracker) PrintStatus() {
	timeDiff := int((time.Since(t.lastTime)).Seconds())
	if !(timeDiff < 1) && t.IsValidWindow(t.lastWindow) {
		fmt.Printf("%vs, %v\n", timeDiff, t.lastWindow)
		t.lastTime = time.Now()
	}
}

func (t *Tracker) IsValidWindow(lastWindow string) (isValid bool) {
	isValid = true

	for _, val := range t.ignorePatterns {
		if str, _ := val.(string); strings.Contains(lastWindow, str) {
			isValid = false
			break
		}
	}

	return
}
