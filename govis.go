package main

import (
	"encoding/json"
	"fmt"
	"github.com/docopt/docopt.go"
	"github.com/op/go-logging"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var log = logging.MustGetLogger("Govis")
var format = logging.MustStringFormatter("%{color}[%{id:03x}] %{time:15:04:05} %{level:.4s}%{color:reset} %{message}")

var configFileName = "govis.json"

func main() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)

	usage := `Govis.
	Usage: Govis [options]

	Options:
		-h --help		Display this help message
	`
	arguments, err := docopt.Parse(usage, nil, true, "Govis 0.0.0", false)

	if err != nil {
		log.Error("Could not parse arguments: %s", err)
	}

	fmt.Println(arguments)

	c := JsonCfg{}
	config := c.GetConfigFile(configFileName)

	TrackTime(config)
}

type JsonCfg struct {
	TickInterval int
	LogName      string
	MinIdleTime  int
}

func (I *JsonCfg) GetConfigFile(path string) *JsonCfg {
	b, err := ioutil.ReadFile(path)

	if err != nil {
		log.Error("%s", err)
	}

	err = json.Unmarshal(b, &I)

	if err != nil {
		log.Error("%s", err)
	}
	return I
}

func TrackTime(config *JsonCfg) {
	minIdleTime, err := time.ParseDuration(strconv.Itoa(config.MinIdleTime) + "s")

	if err != nil {
		log.Error("Could not parse MinIdleTime: err", err)
	}

	lastTime := time.Now()
	lastWindow := GetCurrentWindowName()

	c := time.Tick(time.Duration(config.TickInterval) * time.Millisecond)
	for now := range c {
		currentWindow := GetCurrentWindowName()

		if currentWindow != lastWindow && GetIdleTime() < minIdleTime {
			timeDiff := time.Since(lastTime).String()
			fmt.Printf("Changed from [%s] to [%s] for [%v]\n", lastWindow, currentWindow, timeDiff)
			lastTime = time.Now()
		}

		lastWindow = currentWindow
		var _ = now
	}
}

func GetCurrentWindowID() string {
	out, err := exec.Command("xprop", "-root", "_NET_ACTIVE_WINDOW").Output()

	if err != nil {
		log.Warning("Could get current window id: %s", err)
		return ""
	}

	re := regexp.MustCompile(`0x[a-f0-9]+`)
	match := re.FindStringSubmatch(string(out))

	if len(match) < 1 {
		log.Warning("Could get current window id")
		return ""
	}

	return match[0]
}

func GetCurrentWindowName() string {
	out, err := exec.Command("xprop", "-id", GetCurrentWindowID(), "_NET_WM_NAME").Output()

	if err != nil {
		log.Warning("Could get current window name: %s", err)
		return ""
	}

	re := regexp.MustCompile(`"(.*)"`)
	match := re.FindStringSubmatch(string(out))

	if len(match) < 2 {
		log.Warning("Could get current window name")
		return ""
	}

	return match[1]
}

func GetIdleTime() (t time.Duration) {
	out, err := exec.Command("xprintidle").Output()

	if err != nil {
		log.Warning("Couldn't get idle time: %s", err)
		return
	}

	t, err = time.ParseDuration(strings.TrimSpace(string(out)) + "ms")

	if err != nil {
		log.Warning("Couldn't get idle time: %s", err)
		return
	}

	return
}
