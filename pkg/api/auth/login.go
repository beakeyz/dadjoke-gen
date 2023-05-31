package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/beakeyz/dadjoke-gen/pkg/cookies"
	"github.com/beakeyz/dadjoke-gen/pkg/crypto"
	"github.com/beakeyz/dadjoke-gen/pkg/database"
	"github.com/beakeyz/dadjoke-gen/pkg/settings"
	"github.com/beakeyz/dadjoke-gen/pkg/structs"
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
		//ctx.Resp.Write([]byte("bro u aint got no session fam, fuck off"))
		ctx.JSON(201, nil)
		return
	}

	if !ctx.UserSession.LinkedUser.IsAnonymous {
		//ctx.Resp.Write([]byte("Already logged in!"))
		ctx.JSON(202, nil)
		return
	}

	// hihi checking if we have a sql connection with a single bool, yes my brain works at 2AM I swear
	if database.HasDatabase == false {
		// sql server down
		ctx.JSON(269, nil)
		return
	}

	var attempt LoginAttempt

	json.NewDecoder(ctx.Req.Body).Decode(&attempt)

	var user, origSessId = structs.GetFromDb(attempt.Username)
	if user.IsNull {
		// huilen ;-;
		//ctx.Resp.Write([]byte("That user is not found!"))
		ctx.JSON(203, nil)
		return
	}

	var token string = settings.CONFIG.AuthHash + "" + origSessId.String() + "" + attempt.Password
	token = crypto.HashString(token)
	if user.Token != token {
		fmt.Println(user.Token)
		fmt.Println(token)
		//ctx.Resp.Write([]byte("Password doesn't match!"))
		ctx.JSON(204, nil)
		return
	}

	var session structs.Session = *structs.CreateSessionTemplate(user, uuid.New())

	existingSession, err := structs.SessionRequestGetFromUser(session.LinkedUser)
	// if a valid usersession already existed, use that to log in
	if err == nil && !existingSession.IsNull && !existingSession.LinkedUser.IsAnonymous {
		structs.SessionRequestRemove(*ctx.UserSession)
		cookies.SessionCookie(existingSession.SessionId.String(), ctx.Context, time.Hour*time.Duration(24))
	} else {
		structs.SessionRequestAdd(session)
		structs.SessionRequestRemove(*ctx.UserSession)
		cookies.SessionCookie(session.SessionId.String(), ctx.Context, time.Hour*time.Duration(24))
	}

	//ctx.Resp.Write([]byte("Logged in: " + user.Username))
	ctx.JSON(200, nil)
}
