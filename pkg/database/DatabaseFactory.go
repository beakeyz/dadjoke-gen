package database

import (
	"database/sql"
	"fmt"

	"github.com/beakeyz/dadjoke-gen/pkg/setting"
)

func CraftDatabase(conf *setting.ServerConfig) bool{
  fmt.Println("Trying to craft database...")
  db, initError := sql.Open("mysql", make_connection_string_no_database(conf))
  if initError != nil{
    fmt.Errorf("HAHA LMAO U DED SON")
    return false
  }
	if err := db.Ping(); err != nil{
		return false
	}
  defer db.Close()

  fmt.Println("table not there")
	var err error
	var name string = conf.Sql_database
	_,err = db.Exec("CREATE DATABASE "+name)
  if err != nil {
    //TODO: not ignore
    return false
  }

  _,err = db.Exec("USE "+name)
  if err != nil {
    //TODO: not ignore
    return false
  }

  _,err = db.Exec("CREATE TABLE " + conf.Sql_table + " ( Summary varchar(100), Joke varchar(500), AdditionDate Date, JokeIndex integer )")
  if err != nil {
    //TODO: not ignore
    return false
  }
  return true
}
