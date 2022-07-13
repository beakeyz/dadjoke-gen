package middleware

import (
	"fmt"

	"github.com/beakeyz/dadjoke-gen/pkg/web"
)

func TestMiddleware() web.Handler {
  return func (ctx *web.ReqContext) {
    // TODO make this not stink
    fmt.Println("Hi from the test middleware")
    
  }
}
