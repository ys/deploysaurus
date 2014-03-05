package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/auth"
	"github.com/martini-contrib/binding"
	"github.com/ys/deploysaurus/deploysaurus"
	"net/http"
	"os"
)
import _ "github.com/joho/godotenv/autoload"

func main() {
	m := martini.Classic()
	m.Use(auth.Basic("hooks", os.Getenv("HOOK_KEY")))
	m.Post("/hooks", checkEvent(), binding.Json(deploysaurus.Event{}), binding.ErrorHandler, func(event deploysaurus.Event) string {
		fmt.Println(event.Who(), "deploys", event.What(), ":", event)
		return deploysaurus.Deploy(event)
	})
	http.ListenAndServe(":"+os.Getenv("PORT"), m)
}

func checkEvent() martini.Handler {
	return func(context martini.Context, res http.ResponseWriter, req *http.Request) {
		eventType := req.Header.Get("X-GitHub-Event")
		switch eventType {
		case "deployment":
		default:
			res.WriteHeader(200)
			res.Write([]byte("NOK"))
		}
	}
}
