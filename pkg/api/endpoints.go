package api

import (
	"github.com/beakeyz/dadjoke-gen/pkg/api/auth"
	"github.com/beakeyz/dadjoke-gen/pkg/api/routing"
	"github.com/beakeyz/dadjoke-gen/pkg/web"
)

func (srv *HttpServer) registerApiEndpoints() {
	var rh routing.RouteHandler = srv.RouteHandler

	rh.RegisterGet("/api/test", Test)

	// authentication related routes
	rh.RegisterPost("/api/signup", auth.Register)
	rh.RegisterPost("/api/login", auth.Login)
	rh.RegisterPost("/api/logout", auth.Logout)

	// TODO: protected routes (api functions)

	// TODO: protected routes (admin functions)
}

func Test(ctx *web.ReqContext) {

	if ctx.UserSession.IsNull || ctx.UserSession.LinkedUser.IsNull || ctx.UserSession.LinkedUser.IsAnonymous {
		ctx.Resp.Write([]byte("no"))
	} else {
		ctx.Resp.Write([]byte("yes"))
	}
}
