package notification

type BuildStatusNotification interface {
	BuildSuccess(buildId string)

	BuildFailed(buildId string)
}

