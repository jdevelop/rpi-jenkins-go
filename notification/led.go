package notification

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os"
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
}

func NewPi(piOk PiLed, piFail PiLed) Pi {
	return Pi{piOk, piFail}
}

func (p Pi) BuildSuccess(buildId string) {
	fmt.Println("SUCCESS:" + buildId)
	p.piOkLed.Led.High()
	p.piFailLed.Led.Low()
}

func (p Pi) BuildFailed(buildId string) {
	fmt.Println("FAIL:" + buildId)
	p.piOkLed.Led.Low()
	p.piFailLed.Led.High()
}
