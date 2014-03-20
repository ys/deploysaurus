package deploysaurus

import (
	"encoding/base64"
	"fmt"
	"time"
)

type DbUser struct {
	Id                 string
	Email              string
	GitHubId           string
	GitHubLogin        string
	GitHubToken        string
	HerokuId           string
	HerokuToken        string
	HerokuRefreshToken string
	HerokuExpiration   time.Time
	Authenticated      bool
}

func (u DbUser) GitHubAuthorization() string {
	data := []byte(fmt.Sprintf("%s:x-oauth-basic", u.GitHubToken))
	return base64.StdEncoding.EncodeToString(data)
}

func (u DbUser) HerokuAuthorization() string {
	data := []byte(fmt.Sprintf(":%s", u.HerokuToken))
	return base64.StdEncoding.EncodeToString(data)
}
