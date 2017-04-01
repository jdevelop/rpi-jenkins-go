package main

import (
	bs "github.com/jdevelop/rpi-jenkins-go/buildstatus"
	ntf "github.com/jdevelop/rpi-jenkins-go/notification"
	"math/rand"
	"time"
	"flag"
	"strings"
	"github.com/spf13/viper"
)

func authConfig() (string, string) {
	viper.SetConfigName("auth")
	viper.AddConfigPath("$HOME/.jenkins-pi")
	err := viper.ReadInConfig()
	if err != nil {
		return "", ""
	} else {
		return viper.GetString("pi.jenkins.api.username"), viper.GetString("pi.jenkins.api.accesskey")
	}
}

func displayBuildStatus(status bs.JenkinsBuildStatus, ntf ntf.BuildStatusNotification) {
	switch status.Status {
	case "SUCCESS":
		ntf.BuildSuccess(status)
	case "FAILURE":
		ntf.BuildFailed(status)
	}
}

// -----------------------------------------------------------------------------------------------------

func setup() (bs.BuildStatusProvider, ntf.BuildStatusNotification) {
	urlPtr := flag.String("url", "", "URL for Jenkins")
	piOkPtr := flag.Int("led-success", -1, "Success LED pin number")
	piFailPtr := flag.Int("led-failure", -1, "Failed LED pin number")
	lcdDataPins := flag.String("lcd-data-pin", "", "LCD Data Pins, comma-separated")
	lcdEPin := flag.String("lcd-e-pin", "", "LCD strobe pin")
	lcdRsPin := flag.String("lcd-rs-pin", "", "LCD strobe pin")
	flag.Parse()
	var (
		statusNotifier ntf.BuildStatusNotification
		statusProvider bs.BuildStatusProvider
	)

	if *lcdDataPins != "" {
		statusNotifier = ntf.NewLCD(*lcdRsPin, *lcdEPin, strings.Split(*lcdDataPins,","))
	} else if *piOkPtr == -1 || *piFailPtr == -1 {
		statusNotifier = ntf.NewConsole()
	} else {
		piOk, piFail := ntf.SetupLeds(*piOkPtr, *piFailPtr)
		statusNotifier = ntf.NewPi(piOk, piFail)
	}
	if *urlPtr == "" {
		rand.Seed(time.Now().UTC().UnixNano())
		statusProvider = bs.FakeBuildStatus{}
	} else {
		var username, apikey = authConfig()
		statusProvider = bs.NewJenkinsBuildStatus(*urlPtr, username, apikey)
	}

	return statusProvider, statusNotifier
}

func main() {
	provider, ntfImpl := setup()
	buildStatus := provider.ResolveBuildStatus()
	displayBuildStatus(buildStatus, ntfImpl)
}
