package notification

import (
	"fmt"
	bs "github.com/jdevelop/rpi-jenkins-go/buildstatus"
)

type Console struct{}

func NewConsole() Console {
	return Console{}
}

func (c Console) BuildSuccess(idx int, status bs.JenkinsBuildStatus) {
	fmt.Printf("%d : %s : %d\n", idx, status.Status, status.BuildId)
}

func (c Console) BuildFailed(idx int, status bs.JenkinsBuildStatus) {
	c.BuildSuccess(idx, status)
}
