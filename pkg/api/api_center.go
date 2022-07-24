package api

import (
	"log"
	"os"

	"github.com/beakeyz/dadjoke-gen/pkg/api/auth"
	"github.com/beakeyz/dadjoke-gen/pkg/api/routing"
	"github.com/beakeyz/dadjoke-gen/pkg/middleware"
	"github.com/beakeyz/dadjoke-gen/pkg/web"
)

var tempLog = log.New(os.Stdout,"api_center ", 0)

func (httpServ *HttpServer)RegisterApiEndpoints(){

  tempLog.Println("hi")

  // set the handler
  var rh routing.RouteHandler = httpServ.RouteHandler

  // all the route specific middleware gets defined here
  // NOTE: we still want all the middlewares to live in the middleware package
  test := middleware.TestMiddleware()

  // all the routes get added to the routehandler here
  rh.RegisterGet("/", test, httpServ.Index) 
  rh.RegisterPost("/reg", auth.Register)
  rh.RegisterPost("/login", auth.Login)
  rh.RegisterGet("/protected", ProtectedRoute)
  rh.RegisterGet("/api/v1/get_jokes", httpServ.GetJokes)
  rh.RegisterPost("/api/v1/post_jokes", httpServ.Post_Jokes)

}

func ProtectedRoute (ctx *web.ReqContext) {
  if ctx.UserSession.IsNull {
    ctx.Resp.Write([]byte("wtf u dont have a session????"))
    return
  }

  if ctx.UserSession.LinkedUser.IsAnonymous || ctx.UserSession.LinkedUser.IsNull {
    ctx.Resp.Write([]byte("Can't access this page without an account!"))
    return
  }
  ctx.Resp.Write([]byte("Wow, you reached the secret!!! \nThe secret is: ass == ass at ALL times"))
}
