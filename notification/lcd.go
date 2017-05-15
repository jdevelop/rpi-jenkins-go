package notification

import (
	"fmt"
	l "github.com/jdevelop/golang-rpi-extras/lcd_hd44780"
	bs "github.com/jdevelop/rpi-jenkins-go/buildstatus"
)

type LCD struct {
	l.PiLCD4
}

func NewLCD(rsPin int, enablePin int, data []int) LCD {
	lcd, err := l.NewLCD4(data, rsPin, enablePin)

	if err != nil {
		panic(err.Error())
	}

	lcd.Init()
	return LCD{lcd}
}

func printStatus(r *LCD, status *bs.JenkinsBuildStatus, line int) {
	r.SetCursor(uint8(line), 0)
	r.Print(fmt.Sprintf("BUILD: %d / %s", status.BuildId, status.Status))
}

func (r LCD) BuildSuccess(idx int, status bs.JenkinsBuildStatus) {
	printStatus(&r, &status, idx)
}

func (r LCD) BuildFailed(idx int, status bs.JenkinsBuildStatus) {
	printStatus(&r, &status, idx)
}
