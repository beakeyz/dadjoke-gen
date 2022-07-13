package main

import (
	"os"
	//"github.com/beakeyz/dadjoke-gen/pkg/server/connection"
	//"github.com/beakeyz/dadjoke-gen/pkg/server/routes"
	//"github.com/beakeyz/dadjoke-gen/pkg/server/setup"
	//"github.com/gorilla/mux"
)

func main() {
	/*
	setup.Settings_list.Init()

	connection.SetupConnection()

	if &connection.Connection != nil {
		defer connection.Connection.Close()
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/fetch_joke/{index}", routes.Fetch_joke).Methods("GET")
	r.HandleFunc("/api/fetch_jokes/{amount}", routes.Fetch_jokes).Methods("GET")

	fmt.Println("Server launched")
	http.ListenAndServe(":4000", r)
	*/
	//create server with parameters
	//remove previous code
	os.Exit(Bootstrap())
}
