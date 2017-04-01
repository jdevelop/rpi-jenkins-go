package buildstatus

type BuildStatusProvider interface {
	ResolveBuildStatus() JenkinsBuildStatus
}

type JenkinsBuildStatus struct {
	Status  string
	BuildId int
}
