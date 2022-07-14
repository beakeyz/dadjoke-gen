package middleware

import (
	"fmt"
	"time"

	"github.com/beakeyz/dadjoke-gen/pkg/cookies"
	"github.com/beakeyz/dadjoke-gen/pkg/structures"
	"github.com/beakeyz/dadjoke-gen/pkg/web"
	"github.com/google/uuid"
)

func AuthEntry () web.Handler {
  return func (ctx *web.ReqContext) {

    // TODO

    // check if there already is a session
    // by dehashing the hash that the client sent us
    // and matching it against our local sessionstore

    mngr, createErr := structures.CreateSassManager() 
    if createErr != nil {
      // cry
      fmt.Println("failed to create a new SessionManager!")
      fmt.Println(createErr.Error())
    }

    if !ctx.UserSession.IsNull {
      fmt.Println("Session in the header isn't null")

      sessionId := ctx.UserSession.SessionId

      session, scanErr := mngr.GetSession(&structures.User{
        Username: "Undetermined",
        Token: sessionId,
        IsAnonymous: true,
      })
      if scanErr != nil || session.IsNull == true {
        // cry, we got an invalid sessionId ;-;
        ctx.UserSession = structures.EmptySession()
        fmt.Println("we have recieved an invalid SessionId")
        return
      }

      // Set the right session in the context
      ctx.UserSession = session 

    } else {
      
      usr := &structures.User{
        Username: "Anonymous",
        Token: uuid.New(),
        IsAnonymous: true,
      }

      // cry bc error ='[
      mngr.AddSession(usr)

      cookies.SessionCookie(usr.Token.String(), ctx.Context, time.Hour)
      fmt.Println("Finished assigning a new Session")
      // TODO: check if ctx.Redirect actually works =D
      // ctx.Redirect("/")
    }

    // assign a new session if it didn't already exist

    // make sure the client saves a hash of the sessionId client-side

  }
}

