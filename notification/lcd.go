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
