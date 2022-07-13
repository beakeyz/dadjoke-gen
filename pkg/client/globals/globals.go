package globals

import (
	"time"
	"os"
	"encoding/json"
	"io/ioutil"

	"github.com/beakeyz/dadjoke-gen/pkg/structures"
)

type Token struct{
	TokenValue int `json:"TokenValue"`
	TokenEmail string `json:"TokenEmail"`
	TokenUsername string `json:"TokenUsername"`
}

type Cache struct {
	LastRunDate       string `json:"LastRunDate"`
	ActiveUrl string `json:"CurrentActiveUrl"`
	PreviousJokeIndex int `json:"PreviousJokeIndex"`
	LocalToken* Token `json:"Token"`
	HasToken bool `json:"HasToken"`
	//this might be retarded (most likely) but i'll keep it in anyway because I can't be fucked
	HasSelf bool `json:"HasCache"`
}

const (
	Jokes_path = "assets/jokes.json"
	Cache_path = "assets/cache.json"
)

var (
	Default_cache       = GetDummyCache()
	Previous_joke_index int
	Jokes structures.JokeList
	Glob_Cache Cache
	LocalToken Token
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//Utilities functions
/////////////////////////////////////////////////////////////////////////////////////////////////////////////



/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//Cache functions
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (self *Cache) InitNewCache() Cache{
	date := time.Now()
	formated_date := date.Format("01-02-2006")

	self = &Cache{formated_date, "invalid", -1, &Token{-1, "null", "null"}, false, true}
	return *self
}

func (self *Cache) UpdateCacheDate() Cache {
	date := time.Now()
	formated_date := date.Format("01-02-2006")

	self.LastRunDate = formated_date
	return *self
}

func (self *Cache) UpdateLocalCache() Cache {
	date := time.Now()
	formated_date := date.Format("01-02-2006")

	if self.LocalToken.isValid(){
		self.HasToken = true;
	}else{
		self.HasToken = false;
	}

	self.HasSelf = IsCurrentCacheValid()
	
	cache := Cache{formated_date, "invalid", self.PreviousJokeIndex, self.LocalToken, self.HasToken, self.HasSelf}
	return cache
}

/*
Generates an empty cache lol
*/
func GetDummyCache() Cache{
	var t Token
	x := Cache{"nul", "invalid", -1, t.EmptyToken(), false, false}
	return x
}

/*
Megabrain function (not) that validates the current local cache
*/
func IsCurrentCacheValid() bool{
	var dummy Cache
	file, err := os.Open(Cache_path)

	if (err != nil){
		cache_byte_arr, _ := ioutil.ReadAll(file)

		json.Unmarshal(cache_byte_arr, &dummy)

		if (dummy.LastRunDate != "-1"){
			return true
		}

	}else{
		//TODO: throw error and try to fix the problem
		return false
	}
	return false
}

/*
Validate AND correct the local cache
*/
func (self *Cache) CorrectCache() {
	self.UpdateCacheDate()
	//Still unused	
	self.PreviousJokeIndex = -1

	if !self.LocalToken.isValid() || !self.LocalToken.TokenExists(){
		self.LocalToken = self.LocalToken.EmptyToken()
		self.HasToken = false
	}else{
		self.HasToken = true
	}
}

func (self *Cache) GetUrl() string{
	return self.ActiveUrl	
}


//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//Token functions
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
Inits an empty token
*/
func (self *Token) EmptyToken() *Token {
	t := &Token{-1, "null", "null"}
	self = t;
	return t
}
		
/*
Checks if the token is valid on the client side
*/
func (self Token) isValid() bool {
	return self.TokenValue != -1;
	//TODO: perhaps check with the server if the tokenid is valid, for now this is enough.
}

/*
Checks with the server to see if the users token is valid in the database
*/
func (self *Token) TokenExists() bool {
	//send request and await server response
	//TODO
	return false
}
