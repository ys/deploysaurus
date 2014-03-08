package deploysaurus

import (
	"../heroku"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/common"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/objx"
	"os"
)

var (
	SessionKey string = "USERID"
)

func UseGomniauth() {
	gomniauth.SetSecurityKey(os.Getenv("SECRET_TOKEN"))
	gomniauth.WithProviders(
		github.New(os.Getenv("GITHUB_CLIENT_ID"),
			os.Getenv("GITHUB_CLIENT_SECRET"),
			os.Getenv("GITHUB_REDIRECT_URL")),
		heroku.New(os.Getenv("HEROKU_CLIENT_ID"),
			os.Getenv("HEROKU_CLIENT_SECRET"),
			os.Getenv("HEROKU_REDIRECT_URL")),
	)

}

func GetDistantUser(providerString string, rawQuery string) (common.User, error) {

	provider, err := gomniauth.Provider(providerString)

	if err != nil {
		return nil, err
	}

	oauthParams, err := objx.FromURLQuery(rawQuery)

	if err != nil {
		return nil, err
	}

	creds, err := provider.CompleteAuth(oauthParams)

	if err != nil {
		return nil, err
	}

	return provider.GetUser(creds)
}
