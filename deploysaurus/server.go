package deploysaurus

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/google/go-github/github"
	"github.com/martini-contrib/auth"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/sessions"
	"net/http"
	"os"
)

func LaunchServer(events chan<- Event) {
	m := martini.Classic()
	m.Use(auth.Basic("hooks", os.Getenv("HOOK_KEY")))
	m.Use(sessions.Sessions("_deploysaurus_session", sessions.NewCookieStore([]byte(os.Getenv("SECRET_TOKEN")))))
	m.Use(buildGitHubAuth())
	m.Get("/", oauth2.LoginRequired, handleRoot)
	m.Post("/hooks", checkEvent(), binding.Json(Event{}), binding.ErrorHandler, handleHooks(events))

	http.ListenAndServe(":"+os.Getenv("PORT"), m)
}

func handleRoot(tokens oauth2.Tokens) string {

	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: tokens.Access()},
	}

	client := github.NewClient(t.Client())
	user, _, err := client.Users.Get("")
	if err != nil {
		panic(err)
	}
	return user.String()
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

func buildGitHubAuth() martini.Handler {
	return oauth2.Github(&oauth2.Options{
		ClientId:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"repo_deployment", "repo"},
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL")})
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
