package web

import (
	"github.com/beakeyz/dadjoke-gen/pkg/api/custom_context"
	"github.com/beakeyz/dadjoke-gen/pkg/database"
	"github.com/beakeyz/dadjoke-gen/pkg/logger"
	"github.com/beakeyz/dadjoke-gen/pkg/structs"
	"github.com/google/uuid"
)

var log = *logger.New("Request Handler", 92, false)

// NOTE: shit specific to the contextHandler like different types of handling for different auth methods and requests goes in here
type ContextHandler struct{}

func MakeContextHandler() *ContextHandler {
	return &ContextHandler{}
}

// this mofo is the generally the first thing to interact with the request, so proper scanning is required in case of shit reqz =D
func (ctha *ContextHandler) Middleware(ctx *Context) {

	// Create an COMPLETELY empty context which we will use to build up further down the pipeline
	dummyReqContext := &ReqContext{
		Context:     ctx,
		UserSession: structs.EmptySession(),
		Connection:  &database.Connection,
	}

	// DEBUG
	log.Out("\033[1mRequest with UA \033[0m" + ctx.Req.UserAgent())

	// Try and find a session cookie in the request
	var presumedSessionHashedId string
	cockie := ctx.GetCookie("session")
	if cockie != "" && len(cockie) != 0 {
		log.Out("Found SessionId in header")
		presumedSessionHashedId = cockie
	}

	if sess := ctx.Req.Header.Get("session"); sess != "" {
		presumedSessionHashedId = sess
	}

	var sessionId, err = uuid.Parse(presumedSessionHashedId)
	if err == nil {
		// We found some kind of SessionId. Put it in the context
		// We need to figure out later if this SessionId maps to a user or an Anonymous session
		dummyReqContext.UserSession = &structs.Session{
			SessionId:  sessionId,
			IsNull:     false,
			LinkedUser: *structs.CreateEmptyUser(),
		}
		log.Out("Assigned the sesson id to the context")
	} else {
		log.Out("No session in the request!")
	}

	// pull a funny to the context lmao
	ctx.Req = ctx.Req.WithContext(custom_context.Set(ctx.Req.Context(), dummyReqContext))

}
