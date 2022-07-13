package setting

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//TODO

// DEFAULTS
const(
  SETTINGS = "settings.json"
  ADDRESS = "localhost"
  PORT    = "4000"
  //SQL
)

var (
  CONFIG = ServerConfig{}
  HasInit = false
)

type ServerConfig struct {
  HttpAddress  string `json:"HttpAddress"`
  HttpPort     string `json:"HttpPort"`
  Sql_name     string `json:"SQL_USER"`
	Sql_pass     string `json:"SQL_PASS"`
	Sql_host     string `json:"SQL_HOST"`
	Sql_database string `json:"SQL_DATABASE"`
  Sql_table 	 string `json:"SQL_TABLE"`
  Sql_auth_table string `json:"SQL_AUTH_TABLE"`
  AuthHash     string `json:"AUTH_HASH"`
}

func (self *ServerConfig) LoadFromJson(file string) error {

  stream, err := os.Open(file)
  if err != nil{
    creation_err := self.CreateDefault()
    if creation_err != nil {
      return creation_err
    }
    return err
  }
  raw, _ := ioutil.ReadAll(stream)

  json.Unmarshal(raw, self)
  CONFIG = *self
  HasInit = true
  return nil
}

//TODO
func (self *ServerConfig) CreateDefault() error {
  config := ServerConfig {
    HttpAddress: ADDRESS,
    HttpPort: PORT,
    Sql_name: "<sql login uname>",
    Sql_pass: "<sql login pass>",
    Sql_host: "<sql hostname + port>",
    Sql_database: "<sql main database>",
    Sql_table: "<sql main table>",
    Sql_auth_table: "<sql auth table",
  }
  
  bytes, err := json.Marshal(config)
  if err != nil {
    return err
  }

  write_err := os.WriteFile(SETTINGS, bytes, 0777)
  if write_err != nil {
    return err
  }

  return nil
}
