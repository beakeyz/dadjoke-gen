package database

import (
	"database/sql"

	"github.com/beakeyz/dadjoke-gen/pkg/logger"
	"github.com/beakeyz/dadjoke-gen/pkg/settings"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Connection sql.DB
	//This might be handy in the future, don't judge me
	HasDatabase = false
	mainLog     = *logger.New("SQL Handler", 94, false)
	errLog      = *logger.New("SQL Handler", 94, true)
)

// FIXME: the way this has been done is pure aids. PLEASE move this global Connection variable somewhere more managable lmao
func SetupConnection(conf *settings.ServerConfig) error {
	//NOTE: I KNOW I HAVE DUPLICATE CODE, IM JUST TO LAZY TO CORRECT IT SO GO CRY AB IT
	var err error
	var db *sql.DB
	HasDatabase = false

	db, err = sql.Open("mysql", make_connection_string(conf))
	if err != nil {
		errLog.Out(err.Error())
		return err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		//Check for database error
		errLog.Out(pingErr.Error())
		return pingErr
	}

	//Connect to the database
	Connection = *db

	// To confirm we did actually make a connection
	HasDatabase = true
	mainLog.Out("Connected!")
	return nil
}

// This function gave me brain cancer lmao, I'm just such a smoothbrain
func make_connection_string(conf *settings.ServerConfig) string {
	name := conf.Sql_name
	pass := conf.Sql_pass
	host := conf.Sql_host
	database := conf.Sql_database
	conn := "tcp"
	return name + ":" + pass + "@" + conn + "(" + host + ")/" + database
}

func make_connection_string_no_database(conf *settings.ServerConfig) string {
	name := conf.Sql_name
	pass := conf.Sql_pass
	host := conf.Sql_host
	conn := "tcp"
	return name + ":" + pass + "@" + conn + "(" + host + ")/"
}
