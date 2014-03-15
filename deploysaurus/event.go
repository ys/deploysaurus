package deploysaurus

import (
	"errors"
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

func (event *Event) SenderLogin() string {
	if event.Sender != nil {
		return event.Sender.Login
	} else {
		return "Somebody"
	}
}

func (event *Event) Tarball() string {
	ref := event.Sha
	who, err := event.Who()
	var deployKey string
	if err == nil {
		deployKey = who.GitHubToken
	} else {
		deployKey = ""
	}
	return event.Repository.AuthenticatedArchiveUrl("", ref, deployKey)
}

func (event *Event) What() string {
	if event.Repository == nil {
		return ""
	}
	return event.Repository.FullName
}

func (event *Event) Who() (*DbUser, error) {
	dbUser, err := getUserFromProvider("github", strconv.Itoa(event.Sender.Id))
	return &dbUser, err
}

func (event *Event) Processable() (string, error) {
	sender, err := event.Who()
	if err != nil {
		return "No user for GitHub sender", errors.New("Bad Karma")
	}
	if sender.HerokuId == "" {
		return "User not linked to Heroku, visit http://deploysaurus.yannick.io/auth/heroku when logged in",
			errors.New("Bad karma")
	}
	return "", nil
	//TODO: Verify if app is writable on Heroku for sender Heroku doppelganger
}
