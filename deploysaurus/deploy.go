package deploysaurus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Deployment struct {
	Id        string `json:"id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	Event     *Event
}

type DeployStatus struct {
	Status string `json:"status"`
}

func (d *Deployment) Url() string {
	return fmt.Sprintf("%s/%s", serviceUrl(d.Event.Payload.HerokuApp), d.Id)
}

func (d *Deployment) ResultUrl() string {
	return fmt.Sprintf("%s/result", d.Url())
}

func Deploysaurus(n int, events <-chan Event) {
	for event := range events {
		log.Println("Dinosaur", n, "processing event", event)
		Deploy(event)
	}
}

func Deploy(event Event) string {
	deployment, _ := LaunchDeployment(event)
	ticker := time.NewTicker(time.Second * 20)
	herokuStateToGitHub := map[string]string{"started": "pending",
		"pending":   "pending",
		"succeeded": "success",
		"failed":    "failure",
		"error":     "error"}

	for _ = range ticker.C {
		status := GetDeployState(*deployment, event)
		log.Println(status)
		PostDeploymentStatus(*deployment, event, herokuStateToGitHub[status])
		if herokuStateToGitHub[status] != "pending" {
			ticker.Stop()
		}
	}
	return "OK"
}

func GetDeployState(d Deployment, event Event) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", d.Url(), nil)
	req.Header.Add("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Add("Authorization", authorization(event))
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "error"
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	var ds DeployStatus
	err = json.Unmarshal(contents, &ds)
	if err != nil {
		return "error"
	}
	return ds.Status
}

func PostDeploymentStatus(deployment Deployment, event Event, state string) {
	client := &http.Client{}
	targetUrl := deployment.ResultUrl()
	body := bytes.NewBufferString(fmt.Sprintf(`{"state":"%s", "target_url":"%s"}`, state, targetUrl))
	url := fmt.Sprintf("%s/deployments/%d/statuses", event.Repository.Url, event.Id)
	req, err := http.NewRequest("POST", url, body)
	req.Header.Add("Accept", "application/vnd.github.cannonball-preview+json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", event.Who().GitHubAuthorization()))
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	log.Println(string(contents))
}

func LaunchDeployment(event Event) (*Deployment, error) {
	client := &http.Client{}
	body := bytes.NewBufferString(fmt.Sprintf(`{"source_blob":{"url":"%s"}}`, event.Tarball()))
	req, err := http.NewRequest("POST", serviceUrl(event.Payload.HerokuApp), body)
	req.Header.Add("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Add("Content-Type", "application/json")
	log.Println(authorization(event))
	req.Header.Add("Authorization", authorization(event))
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	var deployment Deployment
	log.Println(string(contents))
	err = json.Unmarshal(contents, &deployment)
	if err != nil {
		return nil, err
	}
	deployment.Event = &event
	return &deployment, nil
}

func serviceUrl(herokuApp string) string {
	return fmt.Sprintf("https://api.heroku.com/apps/%s/builds", herokuApp)
}

func authorization(event Event) string {
	who := event.Who()
	if who == nil {
		log.Println("I can haz a user for Authorization pleazzzz")
	}
	return fmt.Sprintf("Basic %s", who.HerokuAuthorization())
}
