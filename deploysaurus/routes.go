package deploysaurus

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/sessions"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"net/http"
	"strconv"
)

func HandleRoot(user DbUser) (int, interface{}) {
	if user.Authenticated == true {
		return 200, user
	} else {
		return 200, map[string]interface{}{"DaysWithoutAccident": 0, "LastAccident": "Dinosaur attack"}
	}
}

func HandleMe(user DbUser, res http.ResponseWriter) (int, interface{}) {
	return 200, user
}

func HandleLogout(session sessions.Session, res http.ResponseWriter, req *http.Request) {
	session.Delete(SessionKey)
	http.Redirect(res, req, "/", 302)

}
func HandleHooksWrapper(events chan<- Event) martini.Handler {
	return func(event Event, res http.ResponseWriter) (int, interface{}) {
		fmt.Println(event.Who().Email, "want to deploy", event.What(), ":", event)
		message, err := event.Processable()
		if err != nil {
			return 400, map[string]interface{}{"success": false, "message": message}
		}
		events <- event
		return 202, map[string]interface{}{"success": true, "message": "Event dispatched"}
	}
}

func RedirectToProvider(params martini.Params, res http.ResponseWriter, req *http.Request) {
	provider, err := gomniauth.Provider(params["provider"])
	if err != nil {
		panic(err)
	}
	state := gomniauth.NewState("after", "success")
	authUrl, err := provider.GetBeginAuthURL(state, objx.MSI("scope", [1]string{"repo_deployment"}))
	if err != nil {
		panic(err)
	}
	http.Redirect(res, req, authUrl, 302)
}

func CallbackHandler(params martini.Params, req *http.Request, res http.ResponseWriter, s sessions.Session, dbUser DbUser) {
	user, err := GetDistantUser(params["provider"], req.URL.RawQuery)
	if err != nil {
		panic(err)
	}
	creds := user.ProviderCredentials()[params["provider"]]
	if !dbUser.Authenticated {
		var id string
		if params["provider"] == "github" {
			id = strconv.Itoa(int(creds.Get("id").Float64()))
		} else {
			id = creds.Get("id").Str()
		}
		tmpUser, err := GetUserFromProvider(params["provider"], id)
		if err == nil {
			dbUser = tmpUser
		}
	}
	switch params["provider"] {
	case "github":
		dbUser.Email = user.Email()
		dbUser.GitHubToken = creds.Get("access_token").Str()
		dbUser.GitHubId = strconv.Itoa(int(creds.Get("id").Float64()))
		dbUser.GitHubLogin = user.Nickname()
	case "heroku":
		dbUser.Email = user.Email()
		dbUser.HerokuToken = creds.Get("access_token").Str()
		dbUser.HerokuId = creds.Get("id").Str()
		dbUser.HerokuRefreshToken = creds.Get("refresh_token").Str()
	}
	id, _ := SaveUser(dbUser)
	s.Set(SessionKey, id)
	http.Redirect(res, req, "/me", 302)
}
