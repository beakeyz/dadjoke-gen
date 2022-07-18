package api

import (
	"log"
	"os"

	"github.com/beakeyz/dadjoke-gen/pkg/api/auth"
	"github.com/beakeyz/dadjoke-gen/pkg/api/routing"
	"github.com/beakeyz/dadjoke-gen/pkg/middleware"
)

var tempLog = log.New(os.Stdout,"api_center ", 0)

func (httpServ *HttpServer)RegisterApiEndpoints(){

  tempLog.Println("hi")

  // set the handler
  var rh routing.RouteHandler = httpServ.RouteHandler

  // all the route specific middleware gets defined here
  // NOTE: we still want all the middlewares to live in the middleware package
  test := middleware.TestMiddleware()
  authEntry := middleware.AuthEntry()

  // all the routes get added to the routehandler here
  rh.RegisterGet("/", test, httpServ.Index) 
  rh.RegisterGet("/reg", auth.Register)
  rh.RegisterPost("/login", auth.Login)
  rh.RegisterGet("/api/v1/get_jokes", httpServ.GetJokes)
  rh.RegisterPost("/api/v1/post_jokes", httpServ.Post_Jokes)

  rh.RegisterGet("/api/v1/req_ses", authEntry, httpServ.v1_request_session)
}
