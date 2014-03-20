package deploysaurus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Deployment struct {
	Id        string `json:"id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func Deploysaurus(n int, events <-chan Event) {
	for event := range events {
		log.Println("Dinosaur", n, "processing event", event)
		Deploy(event)
	}
}

func Deploy(event Event) string {
	deployment, _ := LaunchDeployment(event)
	PostDeploymentStatus(*deployment, event)
	return "OK"
}

func PostDeploymentStatus(deployment Deployment, event Event) {
	client := &http.Client{}
	targetUrl := fmt.Sprintf("%s/%s/result", serviceUrl(event.Payload.HerokuApp), deployment.Id)
	body := bytes.NewBufferString(fmt.Sprintf(`{"state":"pending", "target_url":"%s"}`, targetUrl))
	url := fmt.Sprintf("%s/deployments/%d/statuses", event.Repository.Url, event.Id)
	log.Println(url)
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
