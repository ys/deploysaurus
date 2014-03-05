package deploysaurus

import (
	"os"
)

type Event struct {
	Id          int         `json:"id"`
	Sha         string      `json:"sha"`
	Name        string      `json:"name"`
	Payload     *Payload    `json:"payload"`
	Description string      `json:"description"`
	Sender      *User       `json:"sender"`
	Repository  *Repository `json:"repository"`
}

type User struct {
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
