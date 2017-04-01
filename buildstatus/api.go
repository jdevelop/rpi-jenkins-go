package buildstatus

import (
	"math/rand"
	"io/ioutil"
	"net/http"
	"github.com/Jeffail/gabs"
	"fmt"
)

type BuildStatusProvider interface {
	ResolveBuildStatus() JenkinsBuildStatus
}

type FakeBuildStatus struct{}

func (st FakeBuildStatus) ResolveBuildStatus() JenkinsBuildStatus {
	diceRoll := rand.Intn(2)
	if diceRoll == 1 {
		return JenkinsBuildStatus{"SUCCESS", 0}
	} else {
		return JenkinsBuildStatus{"FAILURE", 1}
	}
}

// -----------------------------------------------------------------------------------------------------

type JenkinsBuildContext struct {
	url      string
	username string
	password string
}

type JenkinsBuildStatus struct {
	Status  string
	BuildId int
}

func NewJenkinsBuildStatus(urlS string, username string, password string) JenkinsBuildContext {
	return JenkinsBuildContext{urlS, username, password}
}

func (st JenkinsBuildContext) ResolveBuildStatus() JenkinsBuildStatus {
	lastBuild := getLastBuildUrl(retrieveStatus(st, st.url+"/api/json?tree=lastBuild[url]"))
	fmt.Println("Last build URL is " + lastBuild)
	return parseBuildStatus(retrieveStatus(st, lastBuild+"/api/json?tree=result"))
}

func retrieveStatus(conf JenkinsBuildContext, url string) string {
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(conf.username, conf.password)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	respStr, _ := ioutil.ReadAll(resp.Body)
	return string(respStr)
}

func getLastBuildUrl(doc string) string {
	return getParser(doc).Path("lastBuild.url").Data().(string)
}

func parseBuildStatus(doc string) JenkinsBuildStatus {
	parser := getParser(doc)
	status := parser.Path("result").Data().(string)
	return JenkinsBuildStatus{status, 1}
}

func getParser(doc string) (*gabs.Container) {
	parser, err := gabs.ParseJSON([]byte(doc))
	if err != nil {
		panic("Error parsing " + doc + " => " + err.Error())
	}
	return parser
}
