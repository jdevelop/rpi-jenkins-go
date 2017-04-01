package buildstatus

import (
	"math/rand"
)

type FakeBuildStatus struct{}

func (st FakeBuildStatus) ResolveBuildStatus() JenkinsBuildStatus {
	diceRoll := rand.Intn(2)
	if diceRoll == 1 {
		return JenkinsBuildStatus{"SUCCESS", 0}
	} else {
		return JenkinsBuildStatus{"FAILURE", 1}
	}
}
