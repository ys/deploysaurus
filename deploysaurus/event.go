package deploysaurus

import (
	"errors"
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
	Id     int    `json:"id"`
	Login  string `json:"login"`
	DbUser *DbUser
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
	who := event.Who()
	var deployKey string
	if who != nil {
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

func (event *Event) Who() *DbUser {
	if event.Sender != nil {
		return event.Sender.DbUser
	} else {
		return nil
	}
}

func (event *Event) Processable() (string, error) {
	sender := event.Who()
	if sender == nil {
		return "No user for GitHub sender", errors.New("No GitHub")
	}
	if sender.HerokuId == "" {
		return "User not linked to Heroku, visit http://deploysaurus.yannick.io/auth/heroku when logged in",
			errors.New("No Heroku")
	}
	return "", nil
	//TODO: Verify if app is writable on Heroku for sender Heroku doppelganger
}
