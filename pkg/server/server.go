package server

import (
	"context"
	"sync"

	"github.com/beakeyz/dadjoke-gen/pkg/api"
	"github.com/beakeyz/dadjoke-gen/pkg/logger"
	subapps "github.com/beakeyz/dadjoke-gen/pkg/service"
	"github.com/beakeyz/dadjoke-gen/pkg/settings"
	"github.com/beakeyz/dadjoke-gen/pkg/structs"
	"golang.org/x/sync/errgroup"
)

var initErrLog = *logger.New("server-init", 96, true)

type Server struct {
	httpServer     *api.HttpServer
	sessionManager *structs.SessionManager

	sub_services  []subapps.SubService
	ctx           context.Context
	childRoutines *errgroup.Group
	shutdownFn    context.CancelFunc

	mtx     sync.Mutex
	log logger.LogCaller
	errLog  logger.LogCaller
}

type exitCode int

func (srv *Server) Run() exitCode {

	defer srv.shutdown()

	srv.log.Out("Initializing server")

	//setup httpserver
	//   - setup routes
	//   - arm middleware
	//   - serve http protocol
	//
	//create a way to host multiple thingys at once (like grafana did) by wrapping functionality into a service with a predetirmend Run() function

	// funnie limit lol
	srv.childRoutines.SetLimit(2)

	for _, subapp := range srv.sub_services {
		_subapp := subapp
		srv.childRoutines.Go(func() error {

			select {
			case <-srv.ctx.Done():
				return srv.ctx.Err()
			default:
			}

			//srv.rootLog.Out("Running http server on \033[1m" + srv.httpServer.ServerConfig.HttpAddress + ":" + srv.httpServer.ServerConfig.HttpPort + "\033[0m")
			// srv.httpServer.Run(srv.ctx)
			srv.log.Out("Running routine")
			_subapp.Run(srv.ctx)

			srv.log.Out("Stopped the http server")
			return nil
		})
	}

	srv.log.Out("Waiting for the server...")
	if srv.childRoutines.Wait() != nil {
		return -1
	}
	return 1
}

func CreateServer() (error, *Server) {
	rootCtx, shutdownFn := context.WithCancel(context.Background())
	childRoutines, childCtx := errgroup.WithContext(rootCtx)

	var s *Server = &Server{
		httpServer:     nil,
		sessionManager: nil,
		sub_services:   nil,
		ctx:            childCtx,
		childRoutines:  childRoutines,
		shutdownFn:     shutdownFn,
		log:        *logger.New("server", 96, false),
		errLog:         *logger.New("server", 96, true),
	}

	// Load settings
	var servConf *settings.ServerConfig = &settings.ServerConfig{}
	err := servConf.LoadFromJson(settings.SETTINGS)
	if err != nil {
		initErrLog.Out(err.Error())
		return err, s
	}

	// init http server
	HttpServer, httpErr := api.HttpBootstrap(servConf)
	if httpErr != nil {
		initErrLog.Out("Error while Bootstrapping the HttpServer: " + httpErr.Error())
		return err, s
	}

    s.httpServer = HttpServer

	// Setup SessionManager
    // FIXME: having an async session manager is mega dogshit (I think, we could probably have a good impl, but ye this one aint)
	var manager, managerErr = structs.CreateSassManager()
	if managerErr != nil {
		initErrLog.Out(managerErr.Error())
		return err, s 
	}

    s.sessionManager = manager

	subservices := []subapps.SubService{
		HttpServer,
		manager,
	}

    s.sub_services = subservices

	return nil, s
}

func RunServer() exitCode {
    err, mainServer := CreateServer()

    /* Make sure we close the context chanel */
    defer mainServer.shutdownFn()

    /* Failed to create the server, let's dip */
    if (err != nil) {
      return -1
    }

    /* Failed to create a http server, let's dip */
    if (mainServer.httpServer == nil) {
      return -1
    }

    /* Failed to create a session manager, let's dip =/ */
    if (mainServer.sessionManager == nil) {
      return -1
    }

	return mainServer.Run()
}

func (srv *Server) shutdown() error { //TODO: pull funny
	return nil
}
