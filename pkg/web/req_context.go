package web

import (
	"database/sql"
	"strings"

	"github.com/beakeyz/dadjoke-gen/pkg/structures"
)

// TODO: move this mf
// NOTE: this will become usefull when we have intrecate auth systems in place
type ReqContext struct {
  *Context
  UserSession *structures.Session
  Connection *sql.DB
}

func (self *ReqContext) IsApi () bool {
  return strings.HasPrefix(self.Req.URL.Path, "/api")
}
