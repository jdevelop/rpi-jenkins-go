package main

import (
    "github.com/jdevelop/rpi-jenkins-go/rpi"
    "encoding/json"
    "flag"
    "fmt"
    "github.com/spf13/viper"
    "io/ioutil"
    "math/rand"
    "net/http"
    "time"
)

var username, apikey = authConfig()

func retrieveStatus(url string) string {
    client := http.Client{}
    req, _ := http.NewRequest("GET", url, nil)
    req.SetBasicAuth(username, apikey)
    resp, _ := client.Do(req)
    defer resp.Body.Close()
    respStr, _ := ioutil.ReadAll(resp.Body)
    return string(respStr)
}

func asEthernalShit(doc string) map[string]interface{} {
    var jsObj map[string]interface{}
    json.Unmarshal([]byte(doc), &jsObj)
    return jsObj
}

func getLastBuildUrl(doc string) string {
    return asEthernalShit(doc)["lastBuild"].(map[string]interface{})["url"].(string)
}

func parseBuildStatus(doc string) string {
    return asEthernalShit(doc)["result"].(string)
}

func authConfig() (string, string) {
    viper.SetConfigName("auth")
    viper.AddConfigPath("$HOME/.jenkins-pi")
    err := viper.ReadInConfig()
    if err != nil {
        return "", ""
    } else {
        return viper.GetString("pi.jenkins.api.username"), viper.GetString("pi.jenkins.api.accesskey")
    }
}

// -----------------------------------------------------------------------------------------------------

type BuildStatusNotification interface {
    buildSuccess()

    buildFailed()
}

type BuildStatusProvider interface {
    resolveBuildStatus() string
}

// -----------------------------------------------------------------------------------------------------

type FakeBuildStatus struct{}

func (st FakeBuildStatus) resolveBuildStatus() string {
    diceRoll := rand.Intn(2)
    if diceRoll == 0 {
        return "SUCCESS"
    } else {
        return "FAILURE"
    }
}

// -----------------------------------------------------------------------------------------------------

type JenkinsBuildStatus struct {
    url string
}

func (st JenkinsBuildStatus) resolveBuildStatus() string {
    lastBuild := getLastBuildUrl(retrieveStatus(st.url + "/api/json?tree=lastBuild[url]"))
    fmt.Println("Last build URL is " + lastBuild)
    return parseBuildStatus(retrieveStatus(lastBuild + "/api/json?tree=result"))
}

// -----------------------------------------------------------------------------------------------------

type Pi struct {
    piOkLed rpi.PiLed
    piFailLed rpi.PiLed
}

func (p Pi) buildSuccess() {
    fmt.Println("SUCCESS")
    p.piOkLed.Led.High()
    p.piFailLed.Led.Low()
}

func (p Pi) buildFailed() {
    fmt.Println("FAIL")
    p.piOkLed.Led.Low()
    p.piFailLed.Led.High()
}

// -----------------------------------------------------------------------------------------------------

type Console struct{}

func (c Console) buildSuccess() {
    fmt.Println("SUCCESS")
}

func (c Console) buildFailed() {
    fmt.Println("FAILURE")
}

// -----------------------------------------------------------------------------------------------------

func displayBuildStatus(status string, ntf BuildStatusNotification) {
    switch status {
    case "SUCCESS":
        ntf.buildSuccess()
    case "FAILURE":
        ntf.buildFailed()
    }
}

// -----------------------------------------------------------------------------------------------------

func setup() (BuildStatusProvider, BuildStatusNotification) {
    urlPtr := flag.String("url", "", "URL for Jenkins")
    piOkPtr := flag.Int("pi-success", -1, "Success Pin number")
    piFailPtr := flag.Int("pi-failure", -1, "Failed pin number")
    flag.Parse()
    var (
        statusNotificator BuildStatusNotification
        statusProvider    BuildStatusProvider
    )
    if *piOkPtr == -1 || *piFailPtr == -1 {
        statusNotificator = Console{}
    } else {
        piOk, piFail := rpi.SetupLeds(*piOkPtr, *piFailPtr)
        statusNotificator = Pi{ piOkLed: piOk, piFailLed: piFail }
    }
    if *urlPtr == "" {
        statusProvider = FakeBuildStatus{}
    } else {
        statusProvider = JenkinsBuildStatus{url: *urlPtr}
    }
    return statusProvider, statusNotificator
}

func main() {
    rand.Seed(time.Now().UTC().UnixNano())
    provider, ntf := setup()
    buildStatus := provider.resolveBuildStatus()
    displayBuildStatus(buildStatus, ntf)
}
