package deploysaurus

import (
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
