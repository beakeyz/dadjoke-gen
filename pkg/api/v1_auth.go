package api

import (
	"net/http"
)

func (self *HttpServer) AuthRoute(rw http.ResponseWriter, req *http.Request) {
  if req.Header.Get("token") == "" {
    return
  }

  

}
