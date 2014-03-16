package deploysaurus

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/sessions"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"net/http"
	"os"
	"strconv"
)

func HandleRoot(user DbUser, db DB) (int, interface{}) {
	if user.Authenticated == true {
		return 200, user
	} else {
		count, _ := db.GetUsersCount()
		return 200, map[string]interface{}{"DaysWithoutAccident": count,
			"LastAccident": "Dinosaur attack",
			"GitHubAuth":   fmt.Sprintf("%s/auth/github", os.Getenv("DEFAULT_HOST")),
			"HerokuAuth":   fmt.Sprintf("%s/auth/heroku", os.Getenv("DEFAULT_HOST"))}
	}
}

func HandleMe(user DbUser, res http.ResponseWriter) (int, interface{}) {
	return 200, user
}

func HandleLogout(session sessions.Session, res http.ResponseWriter, req *http.Request) {
	session.Delete(SessionKey)
	http.Redirect(res, req, "/", 302)

}
func HandleHooks(db DB, events chan<- Event, event Event, res http.ResponseWriter) (int, interface{}) {
	fmt.Println(event.SenderLogin(), "want to deploy", event.What(), ":", event)
	dbUser, err := db.GetUserFromProvider("github", strconv.Itoa(event.Sender.Id))
	if err != nil {
		return 400, map[string]interface{}{"success": false, "message": "User not found"}
	}
	event.Sender.DbUser = dbUser
	message, err := event.Processable()
	if err != nil {
		return 400, map[string]interface{}{"success": false, "message": message}
	}
	events <- event
	return 202, map[string]interface{}{"success": true, "message": "Event dispatched"}
}

func RedirectToProvider(params martini.Params, res http.ResponseWriter, req *http.Request) {
	provider, err := gomniauth.Provider(params["provider"])
	if err != nil {
		panic(err)
	}
	state := gomniauth.NewState("after", "success")
	moreScopes := map[string]string{"github": "repo, repo_deployment", "heroku": ""}
	authUrl, err := provider.GetBeginAuthURL(state, objx.MSI("scope", moreScopes[params["provider"]]))
	if err != nil {
		panic(err)
	}
	http.Redirect(res, req, authUrl, 302)
}

func CallbackHandler(db DB, params martini.Params, req *http.Request, res http.ResponseWriter, s sessions.Session, dbUser DbUser) {
	user, err := GetDistantUser(params["provider"], req.URL.RawQuery)
	if err != nil {
		http.Redirect(res, req, fmt.Sprintf("/me?%s", err), 302)
	}
	creds := user.ProviderCredentials()[params["provider"]]
	if !dbUser.Authenticated {
		var id string
		if params["provider"] == "github" {
			id = strconv.Itoa(int(creds.Get("id").Float64()))
		} else {
			id = creds.Get("id").Str()
		}
		tmpUser, err := db.GetUserFromProvider(params["provider"], id)
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
	id, _ := db.SaveUser(dbUser)
	s.Set(SessionKey, id)
	http.Redirect(res, req, "/me", 302)
}
