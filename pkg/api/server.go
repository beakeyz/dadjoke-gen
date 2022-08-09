package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/beakeyz/dadjoke-gen/pkg/api/routing"
	"github.com/beakeyz/dadjoke-gen/pkg/database"
	"github.com/beakeyz/dadjoke-gen/pkg/middleware"
	"github.com/beakeyz/dadjoke-gen/pkg/setting"
	"github.com/beakeyz/dadjoke-gen/pkg/structures"
	"github.com/beakeyz/dadjoke-gen/pkg/web"
)

type HttpServer struct {
	mux         *web.Mux
	log         log.Logger
	context     context.Context
	srv         *http.Server
	Middlewares []web.Handler
	Listener    net.Listener

	ContextHandler *web.ContextHandler
	SessionManager *structures.SessionManager
	RouteHandler   routing.RouteHandler
	serverConfig   *setting.ServerConfig
	db             *sql.DB

	done bool
}

func (self *HttpServer) AddMiddleware(handler web.Handler) {
	self.Middlewares = append(self.Middlewares, handler)
}

func BootstrapHttpServer(conf *setting.ServerConfig) (*HttpServer, error) {

	// Setup Database
	err := database.SetupConnection(conf)
	if err != nil {
		fmt.Println("Database error!!")
		fmt.Println(err.Error())
	}

	// Setup SessionManager
	var manager, managerErr = structures.CreateSassManager()
	if managerErr != nil {
		fmt.Errorf(managerErr.Error())
		return &HttpServer{}, managerErr
	}

	// Create HttpServer object
	var serv *HttpServer = &HttpServer{
		mux: web.New(),
		log: *log.New(os.Stdout, "http_server ", 0),
		//wtf
		ContextHandler: web.MakeContextHandler(),
		SessionManager: manager,
		RouteHandler:   routing.CraftRegister(),
		serverConfig:   conf,
		db:             &database.Connection,
	}

	// pull a funny on the routes
	serv.RegisterApiEndpoints()

	// return that badboy
	return serv, nil
}

func (self *HttpServer) Run(ctx context.Context) error {
	self.context = ctx

	fmt.Println("Got to the server! =)")

	self.BootstrapRoutes()

	self.srv = &http.Server{
		Addr:        net.JoinHostPort(self.serverConfig.HttpAddress, self.serverConfig.HttpPort),
		Handler:     self.mux,
		ReadTimeout: time.Duration(time.Millisecond * 1000),
	}

	listener, _ := self.GetListener()
	self.done = true

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		<-ctx.Done()
		if err := self.srv.Shutdown(context.Background()); err != nil {
			self.log.Println("Failed to shut down the server!", err)
		}
	}()

	if err := self.srv.Serve(listener); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			self.log.Println("Exited the server gracefully!")
		}
	}

	wg.Wait()

	self.log.Println("Died =/")
	return nil
}

func (self *HttpServer) WaitForExit() error {

	return nil
}

func (self *HttpServer) GetListener() (net.Listener, error) {
	if self.Listener != nil {
		return self.Listener, nil
	}

	listener, err := net.Listen("tcp", self.srv.Addr)
	if err != nil {
		return nil, fmt.Errorf("Error ocured while setting up the listener %s", err)
	}
	return listener, nil
}

func (self *HttpServer) BootstrapRoutes() {
	fmt.Println("Bootstraped the routes")

	// self.mux.Use(middleware.TestMiddleware())
	self.mux.Use(self.ContextHandler.Middleware)
	self.mux.Use(middleware.AuthEntry())

	for _, mw := range self.Middlewares {
		self.mux.Use(mw)
	}

	self.RouteHandler.ArmRoutes(self.mux)

	fmt.Println("exited bootstrap")
}

func (self *HttpServer) Index(ctx *web.ReqContext) {
	fmt.Println("Index")

	//self.SessionManager.AddSession(&structures.User{
	//  Username: "Anonymous",
	//  Token: uuid.New(),
	//  IsAnonymous: true,
	//})

	//w.Header().Add("token", self.SessionManager.Sessions[0].SessionId.String())
	//fmt.Println(len(self.SessionManager.Sessions))

	if ctx.UserSession.IsNull {
		ctx.Resp.Write([]byte("no Session"))
	} else {
		ctx.Resp.Write([]byte(ctx.UserSession.LinkedUser.Username))
	}
}

func (self *HttpServer) GetJokes(w http.ResponseWriter, rq *http.Request) {

	var jokes structures.JokeList
	w.Header().Set("Content-Type", "application/json")

	rows, err := self.db.Query("SELECT * FROM jokes WHERE 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		jokes.Size++
		if jokes.Size > 10 {
			break
		}
		var tmp structures.Joke
		if err := rows.Scan(&tmp.Summary, &tmp.Joke, &tmp.Date, &tmp.Index); err != nil {
			log.Fatal(err)
		}

		jokes.List = append(jokes.List, tmp)
	}
	json.NewEncoder(w).Encode(jokes)

	fmt.Println("Hit the thing")
}
