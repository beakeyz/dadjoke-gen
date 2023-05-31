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
	rootLog logger.LogCaller
	errLog  logger.LogCaller
}

type exitCode int

func (srv *Server) Run() exitCode {

	defer srv.shutdown()

	srv.rootLog.Out("Initializing...")

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
				srv.rootLog.Out("Im out lmao")
				return srv.ctx.Err()
			default:
			}

			//srv.rootLog.Out("Running http server on \033[1m" + srv.httpServer.ServerConfig.HttpAddress + ":" + srv.httpServer.ServerConfig.HttpPort + "\033[0m")
			// srv.httpServer.Run(srv.ctx)
			srv.rootLog.Out("Running routine")
			_subapp.Run(srv.ctx)

			srv.rootLog.Out("Stopped the http server")
			return nil
		})
	}

	srv.rootLog.Out("Waiting for the server...")
	if srv.childRoutines.Wait() != nil {
		return -1
	}
	return 1
}

func CreateServer() *Server {
	rootCtx, shutdownFn := context.WithCancel(context.Background())
	childRoutines, childCtx := errgroup.WithContext(rootCtx)

	// Load settings
	var servConf *settings.ServerConfig = &settings.ServerConfig{}
	err := servConf.LoadFromJson(settings.SETTINGS)
	if err != nil {
		initErrLog.Out(err.Error())
		// yield
		shutdownFn()
		return &Server{}
	}

	// init http server
	HttpServer, httpErr := api.HttpBootstrap(servConf)
	if httpErr != nil {
		initErrLog.Out("Error while Bootstrapping the HttpServer: " + httpErr.Error())
		// yield
		shutdownFn()
		return &Server{}
	}

	// Setup SessionManager
	var manager, managerErr = structs.CreateSassManager()
	if managerErr != nil {
		initErrLog.Out(managerErr.Error())
		shutdownFn()
		return &Server{}
	}

	subservices := []subapps.SubService{
		HttpServer,
		manager,
	}

	var s *Server = &Server{
		httpServer:     HttpServer,
		sessionManager: manager,
		sub_services:   subservices,
		ctx:            childCtx,
		childRoutines:  childRoutines,
		shutdownFn:     shutdownFn,
		rootLog:        *logger.New("server", 96, false),
		errLog:         *logger.New("server", 96, true),
	}

	return s
}

func RunServer() exitCode {

	var mainServer *Server = CreateServer()
	return mainServer.Run()

}

func (srv *Server) shutdown() error { //TODO: pull funny
	return nil
}
