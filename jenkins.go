package main

import (
	bs "github.com/jdevelop/rpi-jenkins-go/buildstatus"
	ntf "github.com/jdevelop/rpi-jenkins-go/notification"
	"flag"
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

func displayBuildStatus(buildId string, status string, ntf ntf.BuildStatusNotification) {
	switch status {
	case "SUCCESS":
		ntf.BuildSuccess(buildId)
	case "FAILURE":
		ntf.BuildFailed(buildId)
	}
}

// -----------------------------------------------------------------------------------------------------

func setup() (bs.BuildStatusProvider, ntf.BuildStatusNotification) {
	urlPtr := flag.String("url", "", "URL for Jenkins")
	piOkPtr := flag.Int("pi-success", -1, "Success Pin number")
	piFailPtr := flag.Int("pi-failure", -1, "Failed pin number")
	flag.Parse()
	var (
		statusNotificator ntf.BuildStatusNotification
		statusProvider    bs.BuildStatusProvider
	)
	if *piOkPtr == -1 || *piFailPtr == -1 {
		statusNotificator = ntf.NewConsole()
	} else {
		piOk, piFail := ntf.SetupLeds(*piOkPtr, *piFailPtr)
		statusNotificator = ntf.NewPi(piOk, piFail)
	}
	if *urlPtr == "" {
		statusProvider = bs.FakeBuildStatus{}
	} else {
		var username, apikey = authConfig()
		statusProvider = bs.NewJenkinsBuildStatus(*urlPtr, username, apikey)
	}
	return statusProvider, statusNotificator
}

func main() {
	provider, ntfImpl := setup()
	buildId, buildStatus := provider.ResolveBuildStatus()
	displayBuildStatus(buildId, buildStatus, ntfImpl)
}
