package notification

import (
	"fmt"
	bs "github.com/jdevelop/rpi-jenkins-go/buildstatus"
)

type Console struct{}

func NewConsole() Console {
	return Console{}
}

func (c Console) BuildSuccess(status bs.JenkinsBuildStatus) {
	fmt.Printf("%s : %d", status.Status, status.BuildId)
}

func (c Console) BuildFailed(status bs.JenkinsBuildStatus) {
	c.BuildSuccess(status)
}
