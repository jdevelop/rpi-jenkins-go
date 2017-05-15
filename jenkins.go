package main

import (
	bs "github.com/jdevelop/rpi-jenkins-go/buildstatus"
	ntf "github.com/jdevelop/rpi-jenkins-go/notification"
	"math/rand"
	"time"
	"flag"
	"strings"
	"github.com/spf13/viper"
	"strconv"
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

func displayBuildStatus(idx int, status bs.JenkinsBuildStatus, ntf ntf.BuildStatusNotification) {
	switch status.Status {
	case "SUCCESS":
		ntf.BuildSuccess(idx, status)
	case "FAILURE":
		ntf.BuildFailed(idx, status)
	}
}

// -----------------------------------------------------------------------------------------------------

type stringArray []string

type intArray []int

func (arg *stringArray) String() string {
	return ""
}

func (arg *stringArray) Set(src string) error {
	*arg = append(*arg, src)
	return nil
}

func (arg *intArray) String() string {
	return "-1"
}

func (arg *intArray) Set(src string) error {
	res, err := strconv.Atoi(src)
	if err == nil {
		*arg = append(*arg, res)
	}
	return err
}

func setup() ([]bs.BuildStatusProvider, ntf.BuildStatusNotification) {
	var urlPtr stringArray
	flag.Var(&urlPtr, "urls", "URLs for Jenkins")
	var piOkPtr intArray
	flag.Var(&piOkPtr, "led-success", "Success LED pin numbers")
	var piFailPtr intArray
	flag.Var(&piFailPtr, "led-failure", "Failed LED pin numbers")
	lcdDataPins := flag.String("lcd-data-pin", "", "LCD Data Pins, comma-separated")
	lcdEPin := flag.Int("lcd-e-pin", -1, "LCD strobe pin")
	lcdRsPin := flag.Int("lcd-rs-pin", -1, "LCD strobe pin")
	flag.Parse()

	statusNotifier := ntf.NewStack()

	type LcdBuilder func(int)

	var LCD LcdBuilder

	if *lcdDataPins != "" {
		pins := strings.Split(*lcdDataPins, ",")
		intPins := make([]int, len(pins))
		for i, v := range pins {
			intPins[i], _ = strconv.Atoi(v)
		}
		LCD = func(rowNum int) {
			statusNotifier.Register(ntf.NewLCD(*lcdRsPin, *lcdEPin, intPins))
		}
	} else {
		LCD = func(_ int) {}
	}

	if len(piOkPtr) > 0 && len(piFailPtr) > 0 {
		for i, v := range piOkPtr {
			piOk, piFail := ntf.SetupLeds(v, piFailPtr[i])
			statusNotifier.Register(ntf.NewPi(piOk, piFail, i))
		}
	}

	statusNotifier.Register(ntf.NewConsole())

	statusProvider := make([]bs.BuildStatusProvider, 0)

	if len(urlPtr) == 0 {
		rand.Seed(time.Now().UTC().UnixNano())
		statusProvider = append(statusProvider, bs.FakeBuildStatus{})
	} else {
		var username, apikey = authConfig()
		for i, v := range urlPtr {
			statusProvider = append(statusProvider, bs.NewJenkinsBuildStatus(v, username, apikey))
			LCD(i)
		}
	}

	return statusProvider, statusNotifier
}

func main() {
	provider, ntfImpl := setup()
	for i, p := range provider {
		buildStatus := p.ResolveBuildStatus()
		displayBuildStatus(i, buildStatus, ntfImpl)
	}
}
