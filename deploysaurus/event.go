package deploysaurus

import (
	"errors"
	"os"
	"strconv"
)

type Event struct {
	Id          int         `json:"id"`
	Sha         string      `json:"sha"`
	Name        string      `json:"name"`
	Payload     *Payload    `json:"payload"`
	Description string      `json:"description"`
	Sender      *Sender     `json:"sender"`
	Repository  *Repository `json:"repository"`
}

type Sender struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
}

type Payload struct {
	HerokuApp string `json:"heroku_app"`
}

func (event *Event) Tarball() string {
	ref := event.Sha
	deployKey := os.Getenv("GITHUB_DEPLOY_KEY")
	return event.Repository.AuthenticatedArchiveUrl("", ref, deployKey)
}

func (event *Event) What() string {
	if event.Repository == nil {
		return ""
	}
	return event.Repository.FullName
}

func (event *Event) Who() string {
	if event.Sender == nil {
		return ""
	}
	return event.Sender.Login
}

func (event *Event) Processable() (string, error) {
	sender, err := GetUserFromProvider("github", strconv.Itoa(event.Sender.Id))
	if err != nil {
		return "No user for GitHub sender", err
	}
	if sender.HerokuId == "" {
		return "User not linked to Heroku, visit http://deploysaurus.yannick.io/auth/heroku when logged in",
			errors.New("Bad karma")
	}
	return "", nil
	//TODO: Verify if app is writable on Heroku for sender Heroku doppelganger
}
