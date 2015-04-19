package main

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/docopt/docopt.go"
	"github.com/op/go-logging"

	"stathat.com/c/jconfig"
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
	var _ = arguments

	config := jconfig.LoadConfig("/home/ben/.govisrc")

	tracker := Tracker{}
	tracker.Start(config)
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
