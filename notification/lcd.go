package notification

import (
	"fmt"
	l "github.com/jdevelop/gobot-lcd/lcd"
	bs "github.com/jdevelop/rpi-jenkins-go/buildstatus"
)

type LCD struct {
	l.PiLCD4
}

func NewLCD(rsPin string, enablePin string, data []string) LCD {
	var dataPins []l.Pin = make([]l.Pin, len(data))
	for i, v := range data {
		dataPins[i] = l.Pin{v}
	}

	lcd := LCD{
		l.PiLCD4{
			DataPins:  dataPins,
			RsPin:     l.Pin{rsPin},
			EnablePin: l.Pin{enablePin},
		},
	}
	lcd.Init()
	return lcd
}

func printStatus(r *LCD, status *bs.JenkinsBuildStatus) {
	fmt.Println(status.Status)
	r.Cls()
	r.Print(status.Status)
	r.SetCursor(1, 0)
	r.Print(fmt.Sprintf("BUILD: %d", status.BuildId))
}

func (r LCD) BuildSuccess(status bs.JenkinsBuildStatus) {
	printStatus(&r, &status)
}

func (r LCD) BuildFailed(status bs.JenkinsBuildStatus) {
	printStatus(&r, &status)
}
