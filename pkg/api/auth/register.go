package auth

import (
	"encoding/json"
	"fmt"

	"github.com/beakeyz/dadjoke-gen/pkg/crypto"
	"github.com/beakeyz/dadjoke-gen/pkg/setting"
	"github.com/beakeyz/dadjoke-gen/pkg/web"
)

type RegisterAttempt struct {
  Username string
  Password string
}

func Register (ctx *web.ReqContext) {
  // I have the BEST debugging skills I know =)
  fmt.Println("Started Register")
  if ctx.UserSession.IsNull {
    fmt.Println("fuckup one")
    return
  }

  if !ctx.UserSession.LinkedUser.IsNull {
    fmt.Println("fuckup two")
    return
  }
  
  var attempt RegisterAttempt

  json.NewDecoder(ctx.Req.Body).Decode(&attempt)

  // TODO: verify password and username so they are valid and cant cause any sql injection

  // The password we obtain from idk, the reqbody or sm
  var passwrd string = attempt.Password
  // assemble the thing
  var token string = setting.CONFIG.AuthHash + "" + ctx.UserSession.SessionId.String() + "" + passwrd
  // hashing func (SHA1 in this case, idk what the best one is so fuck off)
  token = crypto.HashString(token)

  // ye
  fmt.Println("finished hashing the shit")
  fmt.Printf(token)

  _, err := ctx.Connection.Exec("INSERT INTO `users`(`Username`, `CreationSessId`, `Token`) VALUES ( ?, ?, ? )",
    attempt.Username,
    ctx.UserSession.SessionId.String(),
    token)

  if err != nil {
    fmt.Println(err.Error())
    ctx.Resp.Write([]byte("Whoopsie, registering went wrong =["))
    return
  }

  ctx.Resp.Write([]byte(token))
}
