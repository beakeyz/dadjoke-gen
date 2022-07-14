package web

import (
	"fmt"

	"github.com/beakeyz/dadjoke-gen/pkg/api/custom_context"
	"github.com/beakeyz/dadjoke-gen/pkg/structures"
	"github.com/google/uuid"
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
    UserSession: structures.EmptySession(),
  }

  fmt.Println("Pulling the funny on the context :joy:")

  
  var presumedSessionHashedId string
  cockie := ctx.GetCookie("session")
  if cockie != "" && len(cockie) != 0 {
    presumedSessionHashedId = cockie
  }
  fmt.Println(cockie)

  if sess := ctx.Req.Header.Get("session"); sess != "" {
    presumedSessionHashedId = sess
  }

  var sessionId, err = uuid.Parse(presumedSessionHashedId)
  fmt.Println(sessionId)
  if err == nil {
    dummyReqContext.UserSession = &structures.Session{
      SessionId: sessionId,
      IsNull: false,
      LinkedUser: structures.User{
        IsAnonymous: true,
        Token: sessionId,
        Username: "Anonymous",
      },
    }
  } else {
    fmt.Println("No session in the request!")
  }

  // pull a funny to the context lmao
  ctx.Req = ctx.Req.WithContext(custom_context.Set(ctx.Req.Context(), dummyReqContext))

}
