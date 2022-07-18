package auth

import (
	"encoding/json"

	"github.com/beakeyz/dadjoke-gen/pkg/structures"
	"github.com/beakeyz/dadjoke-gen/pkg/web"
)

// this is going to be the format from now on :clown:
type LoginAttempt struct {
  Username string `json:"Username"`
  Password string `json:"Password"`
 }


func Login (ctx *web.ReqContext) {

  var attempt LoginAttempt
  json.NewDecoder(ctx.Req.Body).Decode(&attempt)

  var user *structures.User = structures.GetFromDb(attempt.Username)
  if user.IsNull {
    // huilen ;-;
    ctx.Resp.Write([]byte("That user is not found!"))
    return
  }

  ctx.Resp.Write([]byte("Found: " + user.Username))
}
