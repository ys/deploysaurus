package deploysaurus

import (
	"encoding/json"
	"github.com/codegangsta/inject"
	"github.com/codegangsta/martini"
	"net/http"
	"reflect"
)

type Response struct {
	Status int
	Body   interface{}
}

func RespondWith() martini.ReturnHandler {
	return func(ctx martini.Context, vals []reflect.Value) {
		rv := ctx.Get(inject.InterfaceOf((*http.ResponseWriter)(nil)))
		res := rv.Interface().(http.ResponseWriter)
		r := &Response{Status: vals[0].Interface().(int), Body: vals[1].Interface()}
		WriteJsonResponse(r, res)
	}
}

func WriteJsonResponse(r *Response, res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(r.Body)
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte(`{"error":"Internal Server Error"}`))
	}
	res.WriteHeader(r.Status)
	res.Write(b)
}
