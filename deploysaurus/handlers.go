package deploysaurus

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/sessions"
	"log"
	"net/http"
)

func CheckEvent() martini.Handler {
	return func(context martini.Context, res http.ResponseWriter, req *http.Request) {
		eventType := req.Header.Get("X-GitHub-Event")
		switch eventType {
		case "deployment":
		case "deployment_status":
			response := &Response{Status: 200, Body: map[string]interface{}{"success": true}}
			WriteJsonResponse(response, res)
		default:
			response := &Response{Status: 400, Body: map[string]interface{}{"success": false, "error_message": "This endpoint only supports deployments events"}}
			WriteJsonResponse(response, res)
		}
	}
}

func SessionUser(db DB) martini.Handler {
	return func(s sessions.Session, c martini.Context) {
		userId := s.Get(SessionKey)
		user := DbUser{Authenticated: false}
		var err error
		if userId != nil {
			user, err = db.GetUser(userId.(string))
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
