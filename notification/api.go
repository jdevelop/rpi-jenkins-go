package notification

import bs "github.com/jdevelop/rpi-jenkins-go/buildstatus"

type BuildStatusNotification interface {
	BuildSuccess(status bs.JenkinsBuildStatus)

	BuildFailed(buildId bs.JenkinsBuildStatus)
}
