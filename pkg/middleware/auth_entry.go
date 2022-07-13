package middleware

import "github.com/beakeyz/dadjoke-gen/pkg/web"

func AuthEntry () web.Handler {
  return func (ctx *web.ReqContext) {

    // TODO

    // check if there already is a session
    // by dehashing the hash that the client sent us
    // and matching it against our local sessionstore

    // assign a new session if it didn't already exist

    // make sure the client saves a hash of the sessionId client-side

  }
}

