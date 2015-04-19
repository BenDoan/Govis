package main

import (
	"fmt"
	"strconv"
	"time"

	"stathat.com/c/jconfig"
)

type Tracker struct {
	lastTime   time.Time
	lastWindow string
}

func (t *Tracker) Start(config *jconfig.Config) {
	minIdleTime, err := time.ParseDuration(strconv.Itoa(config.GetInt("MinIdleTime")) + "s")

	if err != nil {
		log.Error("Could not parse MinIdleTime: err", err)
	}

	t.lastTime = time.Now()
	t.lastWindow = GetCurrentWindowName()

	c := time.Tick(time.Duration(config.GetInt("TickInterval")) * time.Millisecond)
	for now := range c {
		currentWindow := GetCurrentWindowName()

		if currentWindow != t.lastWindow && GetIdleTime() < minIdleTime {
			timeDiff := time.Since(t.lastTime).String()
			fmt.Printf("Changed from [%s] to [%s] for [%v]\n", t.lastWindow, currentWindow, timeDiff)
			t.lastTime = time.Now()
		}

		t.lastWindow = currentWindow
		var _ = now
	}
}
