package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/beakeyz/dadjoke-gen/pkg/setting"
	_ "github.com/go-sql-driver/mysql"
)

var(
	Connection sql.DB
	//This might be handy in the future, don't judge me
	HasDatabase = false
)

func SetupConnection(conf *setting.ServerConfig) error {
	//NOTE: I KNOW I HAVE DUPLICATE CODE, IM JUST TO LAZY TO CORRECT IT SO GO CRY AB IT
	var err error
	var db *sql.DB
	db, err = sql.Open("mysql", make_connection_string(conf))
	if err != nil {
		//Check for database error
		if strings.Contains(err.Error(), "database"){
			//Database is not present
			HasDatabase = false
			CraftDatabase(conf)
			SetupConnection(conf)
			return nil
		}
		log.Fatal(err)
    return err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		//Check for database error
		if strings.Contains(pingErr.Error(), "database"){
			//Database is not present
			HasDatabase = false
			CraftDatabase(conf)
			SetupConnection(conf)
			return nil
		}
		log.Fatal(pingErr)
    return pingErr
	}

	//Connect to the database
	Connection = *db
	fmt.Println("Connected!")
  return nil
}

//This function gave me brain cancer lmao, I'm just such a smoothbrain
func make_connection_string(conf *setting.ServerConfig) string {
	name := conf.Sql_name
	pass := conf.Sql_pass
	host := conf.Sql_host
	database := conf.Sql_database
	conn := "tcp"
	return name + ":" + pass + "@" + conn + "(" + host + ")/" + database
}

func make_connection_string_no_database(conf *setting.ServerConfig) string {
	name := conf.Sql_name
	pass := conf.Sql_pass 
	host := conf.Sql_host 
	conn := "tcp"
	return name + ":" + pass + "@" + conn + "(" + host + ")/"
}
