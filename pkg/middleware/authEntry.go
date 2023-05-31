package middleware

import (
	"fmt"
	"time"

	"github.com/beakeyz/dadjoke-gen/pkg/cookies"
	"github.com/beakeyz/dadjoke-gen/pkg/logger"
	"github.com/beakeyz/dadjoke-gen/pkg/structs"
	"github.com/beakeyz/dadjoke-gen/pkg/web"
	"github.com/google/uuid"
)

var log = *logger.New("Auth", 95, false)
var errLog = *logger.New("Auth", 95, true)

func AuthEntry() web.Handler {
	return func(ctx *web.ReqContext) {

		// {
		// TODO
		// check if there already is a session
		// by dehashing the hash that the client sent us
		// and matching it against our local sessionstore
		// }
		// (FUCK YOU NEVER DEHASH JACKSHIT U INCOMPETENT BALLSACK)

		//mngr, createErr := structs.CreateSassManager()
		//if createErr != nil {
		// cry
		//	errLog.Out("Failed to create a new SessionManager! " + createErr.Error())
		//}

		// TODO: this acao should point to our future domain
		ctx.Resp.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
		// -_-
		ctx.Resp.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		// allow more funnies
		ctx.Resp.Header().Add("Access-Control-Allow-Headers", "X-Requested-With,content-type")
		// make sure our connection is not useless
		ctx.Resp.Header().Add("Access-Control-Allow-Credentials", "true")
		// expose the session header via cors
		ctx.Resp.Header().Add("Access-Control-Expose-Headers", "*")

		if !ctx.UserSession.IsNull {
			log.Out("Session in the header isn't null")

			fmt.Printf("USERSESSION IS %d", ctx.UserSession.SessionId.ID())
			// Grab session from the store (TODO: ckeck if we have enough money)
			session, scanErr := structs.SessionRequestGet(*ctx.UserSession)
			if scanErr != nil {
				// cry, we got an invalid sessionId ;-;
				ctx.UserSession = structs.EmptySession()
				errLog.Out("Something went wrong while refreshing sessionStore: " + scanErr.Error())
				return
			}

			if session.IsNull {
				ctx.UserSession = structs.EmptySession()
				cookies.SessionCookie("", ctx.Context, -1)
				log.Out("Invalid session")
				return
			}

			// Set the right session in the context
			ctx.UserSession = &session
			fmt.Printf("FINANAL PRINTF %d", session.SessionId.ID())
			// update usersession
			cookies.SessionCookie(ctx.UserSession.SessionId.String(), ctx.Context, time.Hour*time.Duration(24))
		} else {

			// first request that a client does ALWAYS creates an Anonymous user. These kinds of users have the IsNull flag, cuz they should not be able to access user-only functions.
			var sessionId uuid.UUID = uuid.New()

			usr := &structs.User{
				Username:    "Anonymous",
				Token:       sessionId.String(),
				IsAnonymous: true,
				IsNull:      true,
			}

			var sess *structs.Session = structs.CreateSessionTemplate(usr, sessionId)

			// cry bc potential error ='[
			//mngr.AddSession(sess)
			structs.SessionRequestAdd(*sess)

			cookies.SessionCookie(usr.Token, ctx.Context, time.Hour*time.Duration(24))
			ctx.UserSession = sess
			log.Out("Finished assigning a new Session")

		}
	}
}
