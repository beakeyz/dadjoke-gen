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

      // Check for the sessionId
      // TODO: if the user is Anonymous, check the default sessionId
      //       otherwise use the Token from the user
      session, scanErr := mngr.GetSession(sessionId)

      if scanErr != nil || session.IsNull == true {
        // cry, we got an invalid sessionId ;-;
        ctx.UserSession = structures.EmptySession()
        cookies.SessionCookie("", ctx.Context, -1)
        fmt.Println("we have recieved an invalid SessionId")
        return
      }

      // TODO : check if the dehashing of the sent ID ties together to hashing an upwrd and a sessionId
      // How? : to have this work we have to store the sessionId that was used to create a specific user
      //        together with its hashed pswrd and other ud, inside a db. When this user is logging in, 
      //        we assume that the SessionId connected to its Anonymous user is different from the one 
      //        that corresponds to the user that is being logged into. When the login is successfull 
      //        (Meaning that the given username/email and password match set hashes in the database), 
      //        we should assign the Correct SessionId to the client, namely the hash between the password 
      //        and the original SessionId used to create the user. THIS is what we would then match against.
      // idk if that makes sense but, message to future me, go fucking make sense out of it, because you have no choice biatch

      // Set the right session in the context
      ctx.UserSession = session 
      // update usersession
      cookies.SessionCookie(ctx.UserSession.SessionId.String(), ctx.Context, time.Hour * time.Duration(24))

    } else {
      
      // first request that a client does ALWAYS creates an Anonymous user. These kinds of users have the IsNull flag, cuz they should not be able to access user-only functions.
      var sessionId uuid.UUID = uuid.New()

      usr := &structures.User{
        Username: "Anonymous",
        Token: sessionId.String(),
        IsAnonymous: true,
        IsNull: true,
      }

      var sess *structures.Session = structures.CreateSessionTemplate(usr, sessionId)

      // cry bc potential error ='[
      mngr.AddSession(sess)

      cookies.SessionCookie(usr.Token, ctx.Context, time.Hour * time.Duration(24))
      fmt.Println("Finished assigning a new Session")
      // TODO: check if ctx.Redirect actually works =D
      // ctx.Redirect("/")
    }
  }
}

