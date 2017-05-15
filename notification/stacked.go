package notification

import (
	bs "github.com/jdevelop/rpi-jenkins-go/buildstatus"
)

type stackedNotifier struct {
	data []BuildStatusNotification
}

func NewStack() (sn stackedNotifier) {
	return
}

func (sn *stackedNotifier) Register(notifierPtr BuildStatusNotification) {
	sn.data = append(sn.data, notifierPtr)
}

func (sn stackedNotifier) BuildSuccess(idx int, status bs.JenkinsBuildStatus) {
	for _, v := range sn.data {
		v.BuildSuccess(idx, status)
	}
}

func (sn stackedNotifier) BuildFailed(idx int, status bs.JenkinsBuildStatus) {
	for _, v := range sn.data {
		v.BuildFailed(idx, status)
	}
}
