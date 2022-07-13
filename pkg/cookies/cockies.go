package cookies

import (
	"net/http"
	"net/url"
	"time"

	"github.com/beakeyz/dadjoke-gen/pkg/web"
)

func TouchCookie (cookieName string, value string, writer http.ResponseWriter, maxLifetime int, path string, secure bool) {

  cookie := http.Cookie {
    Name: cookieName,
    Value: value,
    MaxAge: maxLifetime,
    Path: path,
    Secure: secure,
    HttpOnly: true,
    SameSite: http.SameSiteDefaultMode,
  }
  http.SetCookie(writer, &cookie)
}

func DeleteCookie (name string, writer http.ResponseWriter) {
  TouchCookie(name, "", writer, -1, "/", false)
}

func SessionCookie (value string, context *web.Context, maxLifetime time.Duration) {

  var maxAge int
	if maxLifetime <= 0 {
		maxAge = -1
	} else {
		maxAge = int(maxLifetime.Seconds())
	}

  TouchCookie("session", url.QueryEscape(value), context.Resp, maxAge, "/", true)

}
