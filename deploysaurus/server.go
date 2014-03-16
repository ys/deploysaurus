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
	m.Map(events)
	UseGomniauth()
	m.Map(RespondWith())
	db := mapDatabase(m)
	addSession(m, db)
	addRoutes(m)
	http.ListenAndServe(":"+os.Getenv("PORT"), m)
}

func addSession(m *martini.ClassicMartini, db DB) {
	sessionStore := sessions.NewCookieStore([]byte(os.Getenv("SECRET_TOKEN")))
	m.Use(sessions.Sessions("_deploysaurus_session", sessionStore))
	m.Use(SessionUser(db))
}

func mapDatabase(m *martini.ClassicMartini) DB {
	db, err := GetDB()
	if err != nil {
		panic(err)
	}
	m.MapTo(db, (*DB)(nil))
	return db
}

func addRoutes(m *martini.ClassicMartini) {
	m.Get("/", HandleRoot)
	m.Get("/me", HandleMe)
	m.Get("/logout", HandleLogout)
	m.Post("/hooks", auth.Basic("hooks", os.Getenv("HOOK_KEY")),
		CheckEvent(),
		binding.Json(Event{}),
		binding.ErrorHandler,
		HandleHooks)
	m.Get("/auth/:provider", RedirectToProvider)
	m.Get("/auth/:provider/callback", CallbackHandler)
}
