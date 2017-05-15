package notification

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os"
	"github.com/jdevelop/rpi-jenkins-go/buildstatus"
)

type PiLed struct {
	Led rpio.Pin
}

func SetupLeds(ledOkIdx int, ledFailIdx int) (PiLed, PiLed) {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var ledOk = rpio.Pin(ledOkIdx)
	var ledFail = rpio.Pin(ledFailIdx)
	ledOk.Output()
	ledFail.Output()
	ledOk.Low()
	ledFail.Low()
	return PiLed{ledOk}, PiLed{ledFail}
}

type Pi struct {
	piOkLed   PiLed
	piFailLed PiLed
	buildId   int
}

func NewPi(piOk PiLed, piFail PiLed, idx int) Pi {
	return Pi{piOk, piFail, idx}
}

func debug(status buildstatus.JenkinsBuildStatus) {
	fmt.Printf("%s ⇒ %d", status.Status, status.BuildId)
}

func (p Pi) BuildSuccess(idx int, status buildstatus.JenkinsBuildStatus) {
	if idx != p.buildId {
		return
	}
	debug(status)
	p.piOkLed.Led.High()
	p.piFailLed.Led.Low()
}

func (p Pi) BuildFailed(idx int, status buildstatus.JenkinsBuildStatus) {
	if idx != p.buildId {
		return
	}
	debug(status)
	p.piOkLed.Led.Low()
	p.piFailLed.Led.High()
}
