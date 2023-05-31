package structs

import (
	"fmt"

	"github.com/beakeyz/dadjoke-gen/pkg/database"
	"github.com/google/uuid"
)

type User struct {
	// username so we can easily reference a user
	Username string
	// token that is a combination of their password + their SessionId that was used to create the password
	Token string
	// bananas =D
	IsAnonymous bool
	// where they go?
	IsNull bool
}

func CreateUser(name string, token string) *User {
	return &User{
		Username:    name,
		Token:       token,
		IsAnonymous: false,
		IsNull:      false,
	}
}

func (self *User) AuthenticateUser() {
}

func CreateEmptyUser() *User {
	return &User{IsNull: true}
}

func GetFromDb(username string) (*User, uuid.UUID) {

	thing := database.Connection.QueryRow("SELECT * FROM users WHERE Username = ?", username)
	if thing.Err() != nil {
		return CreateEmptyUser(), uuid.UUID{}
	}

	// TODO: fix database entries + thing.Scan fixen
	//var userid int
	var uname string
	var email string
	var token string
	var _origSessId string
	//var usersets string

	if scanErr := thing.Scan(&uname, &email, &_origSessId, &token); scanErr != nil {
		fmt.Println(scanErr.Error())
		return CreateEmptyUser(), uuid.UUID{}
	}

	origSessId, parseErr2 := uuid.Parse(_origSessId)
	if parseErr2 != nil {
		fmt.Println(parseErr2.Error())
		return CreateEmptyUser(), uuid.UUID{}
	}

	var dummyUser User = User{
		Username:    uname,
		Token:       token,
		IsAnonymous: false,
		IsNull:      false,
	}

	return &dummyUser, origSessId
}
