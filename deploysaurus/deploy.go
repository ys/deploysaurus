package deploysaurus

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Deploysaurus(n int, events <-chan Event) {
	for event := range events {
		log.Println("Dinosaur", n, "processing event", event)
		Deploy(event)
	}
}

func Deploy(event Event) string {
	client := &http.Client{}
	body := bytes.NewBufferString(fmt.Sprintf(`{"source_blob":{"url":"%s"}}`, event.Tarball()))
	req, err := http.NewRequest("POST", serviceUrl(event.Payload.HerokuApp), body)
	req.Header.Add("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", authorization(event))
	log.Println(req)
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(contents))
	return "OK"
}

func serviceUrl(herokuApp string) string {
	return fmt.Sprintf("https://api.heroku.com/apps/%s/builds", herokuApp)
}

func authorization(event Event) string {
	data := []byte(fmt.Sprintf(":%s", event.Who().HerokuToken))
	return fmt.Sprintf("Authorization %s", base64.StdEncoding.EncodeToString(data))
}
