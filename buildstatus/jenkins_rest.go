package buildstatus

import (
	"io/ioutil"
	"net/http"
	"github.com/Jeffail/gabs"
	"fmt"
)

type JenkinsBuildContext struct {
	url      string
	username string
	password string
}

func NewJenkinsBuildStatus(urlS string, username string, password string) JenkinsBuildContext {
	return JenkinsBuildContext{urlS, username, password}
}

func (st JenkinsBuildContext) ResolveBuildStatus() JenkinsBuildStatus {
	lastBuild := getLastBuildUrl(retrieveStatus(st, st.url+"/api/json?tree=lastBuild[url]"))
	fmt.Println("Last build URL is " + lastBuild)
	return parseBuildStatus(retrieveStatus(st, lastBuild+"/api/json?tree=result,number"))
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

func parseBuildStatus(doc string) (status JenkinsBuildStatus) {
	parser := getParser(doc)
	status.Status = parser.Path("result").Data().(string)
	status.BuildId = int(parser.Path("number").Data().(float64))
	return
}

func getParser(doc string) (*gabs.Container) {
	parser, err := gabs.ParseJSON([]byte(doc))
	if err != nil {
		panic("Error parsing " + doc + " => " + err.Error())
	}
	return parser
}

