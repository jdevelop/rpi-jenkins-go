package rpi

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
