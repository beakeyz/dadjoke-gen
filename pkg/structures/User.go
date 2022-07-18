package structures

import (
	"fmt"

	"github.com/beakeyz/dadjoke-gen/pkg/database"
	"github.com/google/uuid"
)

type User struct{
  // username so we can easily reference a user
	Username string
  // token that is a combination of their password + their SessionId that was used to create the password
	Token uuid.UUID 
  // bananas =D
  IsAnonymous bool
  // where they go?
  IsNull bool
}

func CreateUser (name string, token uuid.UUID) *User {
  return &User{
    Username: name,
    Token: token,
    IsAnonymous: false,
    IsNull: false,
  }
}

func (self *User) AuthenticateUser () {
}

func CreateEmptyUser () *User {
  return &User{IsNull: true}
}

func GetFromDb (username string) *User {
  
  thing := database.Connection.QueryRow("SELECT * FROM users WHERE Username = ?", username)
  if thing.Err() != nil {
    return CreateEmptyUser()
  }

  // TODO: fix database entries + thing.Scan fixen
  var dummyUser User = *CreateEmptyUser()
  if scanErr := thing.Scan(&dummyUser); scanErr != nil {
    fmt.Println(scanErr.Error())
    return CreateEmptyUser()
  }
  dummyUser.IsNull = false

  return &dummyUser
}
