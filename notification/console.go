package notification

import (
	"fmt"
	"time"
	"math/rand"
)

type Console struct{}

func NewConsole() Console {
	rand.Seed(time.Now().UTC().UnixNano())
	return Console{}
}

func (c Console) BuildSuccess(buildId string) {
	fmt.Println("SUCCESS:" + buildId)
}

func (c Console) BuildFailed(buildId string) {
	fmt.Println("FAILURE: " + buildId)
}
