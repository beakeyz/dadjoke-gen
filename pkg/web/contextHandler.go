package web

import (
	"fmt"

	"github.com/beakeyz/dadjoke-gen/pkg/api/custom_context"
	"github.com/beakeyz/dadjoke-gen/pkg/structures"
)

// NOTE: shit specific to the contextHandler like different types of handling for different auth methods and requests goes in here
type ContextHandler struct {
}

func MakeContextHandler () *ContextHandler {
  return &ContextHandler{}
}

func (self *ContextHandler) Middleware (ctx *Context) {

  dummyReqContext := &ReqContext{
    Context: ctx,
    UserSession: &structures.Session{},
  }

  fmt.Println("Pulling the funny on the context :joy:")

  /*
  var isUsingCustomClient bool
  var presumedSessionHashedId string
  if cockie := ctx.GetCookie("session"); cockie != "" && len(cockie) != 0 {
    presumedSessionHashedId = cockie
    isUsingCustomClient = false
  }

  if sess := ctx.Req.Header.Get("session"); sess != "" {
    isUsingCustomClient = true
    presumedSessionHashedId = sess
  }
  */

  // pull a funny to the context lmao
  ctx.Req = ctx.Req.WithContext(custom_context.Set(ctx.Req.Context(), dummyReqContext))

}
