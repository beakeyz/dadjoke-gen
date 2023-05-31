package web

import (
	"database/sql"
	"strings"

	"github.com/beakeyz/dadjoke-gen/pkg/structs"
)

// TODO: move this mf
// NOTE: this will become usefull when we have intrecate auth systems in place
type ReqContext struct {
	*Context
	UserSession *structs.Session
	Connection  *sql.DB
}

func (reqCtx *ReqContext) IsApi() bool {
	return strings.HasPrefix(reqCtx.Req.URL.Path, "/api")
}
