package middleware

import (
	"fmt"

	"github.com/beakeyz/dadjoke-gen/pkg/cookies"
	"github.com/beakeyz/dadjoke-gen/pkg/web"
)

func AuthEntry () web.Handler {
  return func (ctx *web.ReqContext) {

    // TODO

    // check if there already is a session
    // by dehashing the hash that the client sent us
    // and matching it against our local sessionstore
    if !ctx.UserSession.IsNull {
      fmt.Println("thing")
    } else {
      cookies.TouchCookie("session", "27ecab53-bac5-4864-9b63-70f462b022bf", ctx.Resp, 100000000, "/", true)
      fmt.Println("Cookie thing done =)")
      ctx.Redirect("/")
    }

    // assign a new session if it didn't already exist

    // make sure the client saves a hash of the sessionId client-side

  }
}

