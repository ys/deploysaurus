package deploysaurus

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/auth"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/sessions"
	"net/http"
	"os"
)

func LaunchServer(events chan<- Event) {
	m := martini.Classic()
	sessionStore := sessions.NewCookieStore([]byte(os.Getenv("SECRET_TOKEN")))
	UseGomniauth()
	m.Map(RespondWith())
	m.Use(sessions.Sessions("_deploysaurus_session", sessionStore))
	m.Use(SessionUser())
	m.Get("/", HandleRoot)
	m.Get("/me", HandleMe)
	m.Get("/logout", HandleLogout)
	m.Post("/hooks", auth.Basic("hooks", os.Getenv("HOOK_KEY")), CheckEvent(), binding.Json(Event{}), binding.ErrorHandler, HandleHooksWrapper(events))
	m.Get("/auth/:provider", RedirectToProvider)
	m.Get("/auth/:provider/callback", CallbackHandler)

	http.ListenAndServe(":"+os.Getenv("PORT"), m)
}
