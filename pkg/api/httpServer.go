package api

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/beakeyz/dadjoke-gen/pkg/api/routing"
	"github.com/beakeyz/dadjoke-gen/pkg/bus"
	"github.com/beakeyz/dadjoke-gen/pkg/database"
	"github.com/beakeyz/dadjoke-gen/pkg/logger"
	"github.com/beakeyz/dadjoke-gen/pkg/middleware"
	"github.com/beakeyz/dadjoke-gen/pkg/middleware/eventhandlers"
	"github.com/beakeyz/dadjoke-gen/pkg/settings"
	"github.com/beakeyz/dadjoke-gen/pkg/web"
)

var initErrLog = *logger.New("server-init", 96, true)
var initLog = *logger.New("server-init", 96, false)

// http server structure layout
type HttpServer struct {
	mux    *web.Mux
	log    logger.LogCaller
	errLog logger.LogCaller

	context        context.Context
	internalServer *http.Server
	Listener       net.Listener

	ContextHandler *web.ContextHandler
	RouteHandler   routing.RouteHandler
	ServerConfig   *settings.ServerConfig
	db             *sql.DB
}

// init http server
func HttpBootstrap(conf *settings.ServerConfig) (*HttpServer, error) {

	dbErr := database.SetupConnection(conf)
	if dbErr != nil {
		initLog.Out("Database error, see output above")
	}

	// Create HttpServer object
	var serv *HttpServer = &HttpServer{
		mux:    web.New(),
		log:    *logger.New("http_server", 94, false),
		errLog: *logger.New("http_server", 94, true),
		//wtf
		ContextHandler: web.MakeContextHandler(),
		RouteHandler:   routing.CraftRegister(),
		ServerConfig:   conf,
		db:             &database.Connection,
	}

	// pull a funny on the routes
	serv.registerApiEndpoints()

	serv.initEventHandlers()

	serv.log.Out("HttpServer initialized")
	// return that badboy
	return serv, nil
}

// Run the http server (duh)
func (srv *HttpServer) Run(ctx context.Context) error {
	srv.context = ctx

	srv.initRoutes()

	srv.internalServer = &http.Server{
		Addr:        net.JoinHostPort(srv.ServerConfig.HttpAddress, srv.ServerConfig.HttpPort),
		Handler:     srv.mux,
		ReadTimeout: time.Duration(time.Millisecond * 1000),
	}

	listener, listenerErr := srv.GetListener()
	if listenerErr != nil {
		srv.errLog.Out("Listener error!")
	}

	var dummyGroup sync.WaitGroup
	dummyGroup.Add(1)

	go func() {
		defer dummyGroup.Done()

		<-ctx.Done()
		if shutdownErr := srv.internalServer.Shutdown(context.Background()); shutdownErr != nil {
			srv.log.Out("Failed to shut down the server: " + shutdownErr.Error())
		}
	}()

	if serveErr := srv.internalServer.Serve(listener); serveErr != nil {
		if errors.Is(serveErr, http.ErrServerClosed) {
			srv.log.Out("Exited the server gracefully!")
		} else {
			srv.log.Out("Bro idk what happend but here figure it out: " + serveErr.Error())
		}
	}

	dummyGroup.Wait()
	return nil
}

// create a quick tcp listener
// TODO: look into protocols r sm
func (srv *HttpServer) GetListener() (net.Listener, error) {
	if srv.Listener != nil {
		return srv.Listener, nil
	}

	listener, err := net.Listen("tcp", srv.internalServer.Addr)
	if err != nil {
		srv.log.Out("Error ocured while setting up the listener: " + err.Error())
		return nil, err
	}
	return listener, nil
}

// This is where all the routes get initialized
func (srv *HttpServer) initRoutes() {
	var m *web.Mux = srv.mux

	// Used to create a ReqContext for every request
	m.Use(srv.ContextHandler.Middleware)

	// Used to assign sessions to requests
	m.Use(middleware.AuthEntry())

	// routes
	srv.RouteHandler.ArmRoutes(m)
}

func (srv *HttpServer) initEventHandlers() {
	bus.EVENTBUS.Register(eventhandlers.SimpleUserChecker)
}
