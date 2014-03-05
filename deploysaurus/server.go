package deploysaurus

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/auth"
	"github.com/martini-contrib/binding"
	"net/http"
	"os"
)

func LaunchServer(events chan<- Event) {
	m := martini.Classic()
	m.Use(auth.Basic("hooks", os.Getenv("HOOK_KEY")))
	m.Post("/hooks", checkEvent(), binding.Json(Event{}), binding.ErrorHandler, func(event Event, res http.ResponseWriter) {
		fmt.Println(event.Who(), "deploys", event.What(), ":", event)
		events <- event
		res.WriteHeader(202)
		res.Header().Set("Content-Type", "application/json")
		fmt.Fprint(res, Response{"success": true,
			"message": "Event dispatched"})
		return

	})
	http.ListenAndServe(":"+os.Getenv("PORT"), m)
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
