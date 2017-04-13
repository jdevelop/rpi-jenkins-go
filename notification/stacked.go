package notification

import (
	bs "github.com/jdevelop/rpi-jenkins-go/buildstatus"
)

type stackedNotifier struct {
	data []BuildStatusNotification
	idx  int
}

func NewStack(num int) (sn stackedNotifier) {
	sn.data = make([]BuildStatusNotification, num)
	sn.idx = 0
	return
}

func (sn *stackedNotifier) Register(notifierPtr BuildStatusNotification) {
	sn.data[sn.idx] = notifierPtr
	sn.idx++
}

func (sn stackedNotifier) BuildSuccess(status bs.JenkinsBuildStatus) {
	for i := 0; i < sn.idx; i++ {
		sn.data[i].BuildSuccess(status)
	}
}

func (sn stackedNotifier) BuildFailed(status bs.JenkinsBuildStatus) {
	for i := 0; i < sn.idx; i++ {
		sn.data[i].BuildFailed(status)
	}
}
