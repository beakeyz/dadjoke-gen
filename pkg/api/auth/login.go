package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/beakeyz/dadjoke-gen/pkg/cookies"
	"github.com/beakeyz/dadjoke-gen/pkg/crypto"
	"github.com/beakeyz/dadjoke-gen/pkg/setting"
	"github.com/beakeyz/dadjoke-gen/pkg/structures"
	"github.com/beakeyz/dadjoke-gen/pkg/web"
	"github.com/google/uuid"
)

// this is going to be the format from now on :clown:
type LoginAttempt struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func Login(ctx *web.ReqContext) {

	if ctx.UserSession.IsNull {
		ctx.Resp.Write([]byte("bro u aint got no session fam, fuck off"))
		return
	}

	if !ctx.UserSession.LinkedUser.IsAnonymous {
		ctx.Resp.Write([]byte("Already logged in!"))
		return
	}

	var attempt LoginAttempt

	json.NewDecoder(ctx.Req.Body).Decode(&attempt)

	var user, origSessId = structures.GetFromDb(attempt.Username)
	if user.IsNull {
		// huilen ;-;
		ctx.Resp.Write([]byte("That user is not found!"))
		return
	}

	var token string = setting.CONFIG.AuthHash + "" + origSessId.String() + "" + attempt.Password
	token = crypto.HashString(token)
	if user.Token != token {
		fmt.Println(user.Token)
		fmt.Println(token)
		ctx.Resp.Write([]byte("Password doesn't match!"))
		return
	}

	var session structures.Session = *structures.CreateSessionTemplate(user, uuid.New())

	// create a sessionManager
	mngr, _ := structures.CreateSassManager()

	// add the new session
	mngr.AddSession(&session)
	// remove the old session
	mngr.RemoveSession(ctx.UserSession)

	// send the cookie to the client to confirm authentication
	cookies.SessionCookie(session.SessionId.String(), ctx.Context, time.Hour*time.Duration(24))

	ctx.Resp.Write([]byte("Found: " + user.Username))
}
