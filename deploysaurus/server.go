package deploysaurus

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/auth"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/sessions"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"log"
	"net/http"
	"os"
	"strconv"
)

func LaunchServer(events chan<- Event) {
	m := martini.Classic()
	sessionStore := sessions.NewCookieStore([]byte(os.Getenv("SECRET_TOKEN")))
	UseGomniauth()
	m.Use(sessions.Sessions("_deploysaurus_session", sessionStore))
	m.Use(SessionUser())
	m.Get("/", handleRoot)
	m.Post("/hooks", auth.Basic("hooks", os.Getenv("HOOK_KEY")), checkEvent(), binding.Json(Event{}), binding.ErrorHandler, handleHooks(events))
	m.Get("/auth/:provider", redirectToProvider)
	m.Get("/auth/:provider/callback", callbackHandler)

	http.ListenAndServe(":"+os.Getenv("PORT"), m)
}

func handleRoot(user DbUser) string {
	if user.Authenticated == true {
		return fmt.Sprintf("Hello %", user.GitHubLogin)
	} else {
		return "GOGOGO"
	}
}

func handleHooks(events chan<- Event) martini.Handler {
	return func(event Event, res http.ResponseWriter) {
		fmt.Println(event.Who(), "deploys", event.What(), ":", event)
		events <- event
		res.WriteHeader(202)
		res.Header().Set("Content-Type", "application/json")
		fmt.Fprint(res, Response{"success": true,
			"message": "Event dispatched"})
		return
	}
}

func checkEvent() martini.Handler {
	return func(context martini.Context, res http.ResponseWriter, req *http.Request) {
		eventType := req.Header.Get("X-GitHub-Event")
		if eventType != "deployment" {
			res.WriteHeader(400)
			res.Header().Set("Content-Type", "application/json")
			fmt.Fprint(res, Response{"success": false,
				"error_message": "This endpoint only supports deployments events"})
			return
		}
	}
}
func SessionUser() martini.Handler {
	return func(s sessions.Session, c martini.Context) {
		userId := s.Get(SessionKey)

		user := DbUser{Authenticated: false}
		var err error
		if userId != nil {
			user, err = GetUser(userId.(string))
			if user.Id != "" {
				user.Authenticated = true
			}
			if err != nil {
				log.Printf("Login Error: %v\n", err)
			}
		}

		c.Map(user)
	}
}

func redirectToProvider(params martini.Params, res http.ResponseWriter, req *http.Request) {
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

func callbackHandler(params martini.Params, req *http.Request, s sessions.Session, dbUser DbUser) string {
	user, err := GetDistantUser(params["provider"], req.URL.RawQuery)
	if err != nil {
		panic(err)
	}
	creds := user.ProviderCredentials()[params["provider"]]
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
	return id
}
