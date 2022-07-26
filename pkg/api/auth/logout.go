package auth

import (
	"github.com/beakeyz/dadjoke-gen/pkg/cookies"
	"github.com/beakeyz/dadjoke-gen/pkg/structures"
	"github.com/beakeyz/dadjoke-gen/pkg/web"
)

func Logout (ctx *web.ReqContext)  {
  if ctx.UserSession.IsNull || ctx.UserSession.LinkedUser.IsAnonymous {
    ctx.Resp.Write([]byte("Cant logout if ur not logged in lmaooo"))
    return
  }
  
  cookies.SessionCookie("", ctx.Context, -1)
  var mng, err = structures.CreateSassManager()
  if err != nil {
    ctx.Resp.Write([]byte("Some issues occured while trying to logout =/"))
  }
  mng.RemoveSession(ctx.UserSession)
  ctx.Resp.Write([]byte("Succesfully logged out!"))
}
