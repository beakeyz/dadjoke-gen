package server

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/beakeyz/dadjoke-gen/pkg/api"
	"golang.org/x/sync/errgroup"
)


type Server struct {
  HTTPServer *api.HttpServer

  ctx context.Context
  childRoutines *errgroup.Group
  shutdownFn context.CancelFunc

  mtx sync.Mutex
  log log.Logger
}

func New(HttpServer *api.HttpServer) (*Server, error){
  s, err := newServ(HttpServer)
  if err == nil{
    if _err := s.init(); _err != nil{
      return nil, _err
    }
    return s, nil
  }
  return nil, err
}

func newServ(HttpServer *api.HttpServer) (*Server, error){
  rootCtx, shutdownFn := context.WithCancel(context.Background())
  childRoutines, childCtx := errgroup.WithContext(rootCtx)

  var s *Server = &Server{
    HTTPServer:     HttpServer,
    ctx:            childCtx,
    childRoutines:  childRoutines,
    shutdownFn:     shutdownFn,
    log:            *log.New(os.Stdout, "dadjoke-server ", 0),
  }
  return s, nil
}

//initialize the server (with context / params?)
func (self *Server)init() error {
  self.mtx.Lock()
  defer self.mtx.Unlock()

  //put any initialization functions here that might require error handling
    return nil
}

func (self *Server)Run() error{
  defer self.shutdown()

  self.log.Println("initializing...")
  if e := self.init(); e != nil{
    //should not fail yet
    self.log.Panicln("Failed to init")
    return e
  }

  //setup httpserver
  //   - setup routes
  //   - arm middleware
  //   - serve http protocol
  //
  //create a way to host multiple thingys at once (like grafana did) by wrapping functionality into a service with a predetirmend Run() function

  self.childRoutines.Go(func() error {

    select{
    case <-self.ctx.Done():
      self.log.Println(" Im out lmao")
      return self.ctx.Err()
    default:
    }

    self.log.Println("Running http server!")
    self.HTTPServer.Run(self.ctx)

    self.log.Println("Stopped the http server")
    return nil
  })

  self.log.Println("Waiting for the server...")
  return self.childRoutines.Wait()
}

func (self *Server)shutdown() error{ //TODO: pull funny
  return nil
}
