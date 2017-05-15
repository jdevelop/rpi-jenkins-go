package notification

import bs "github.com/jdevelop/rpi-jenkins-go/buildstatus"

type BuildStatusNotification interface {
	BuildSuccess(idx int, status bs.JenkinsBuildStatus)

	BuildFailed(idx int, buildId bs.JenkinsBuildStatus)
}
