package buildstatus

import (
	"math/rand"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"fmt"
)

type BuildStatusProvider interface {
	ResolveBuildStatus() (string, string)
}

type FakeBuildStatus struct{}

func (st FakeBuildStatus) ResolveBuildStatus() (string, string) {
	diceRoll := rand.Intn(2)
	if diceRoll == 0 {
		return "0", "SUCCESS"
	} else {
		return "1", "FAILURE"
	}
}

// -----------------------------------------------------------------------------------------------------

type JenkinsBuildStatus struct {
	url      string
	username string
	password string
}

func NewJenkinsBuildStatus(urlS string, username string, password string) JenkinsBuildStatus {
	return JenkinsBuildStatus{urlS, username, password}
}

func (st JenkinsBuildStatus) ResolveBuildStatus() (string, string) {
	lastBuild := getLastBuildUrl(retrieveStatus(st, st.url + "/api/json?tree=lastBuild[url]"))
	fmt.Println("Last build URL is " + lastBuild)
	return "1", parseBuildStatus(retrieveStatus(st, lastBuild + "/api/json?tree=result"))
}

func retrieveStatus(conf JenkinsBuildStatus, url string) string {
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(conf.username, conf.password)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	respStr, _ := ioutil.ReadAll(resp.Body)
	return string(respStr)
}

func getLastBuildUrl(doc string) string {
	return asEthernalShit(doc)["lastBuild"].(map[string]interface{})["url"].(string)
}

func parseBuildStatus(doc string) string {
	return asEthernalShit(doc)["result"].(string)
}

func asEthernalShit(doc string) map[string]interface{} {
	var jsObj map[string]interface{}
	json.Unmarshal([]byte(doc), &jsObj)
	return jsObj
}
