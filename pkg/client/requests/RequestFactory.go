package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/beakeyz/dadjoke-gen/pkg/client/globals"
	"github.com/beakeyz/dadjoke-gen/pkg/structures"
)

//This is just so we can queue up a large amount of requests and have a somewhat constant delay between them, while keeping programflow
//I'll need to add some functions to regulate this queue later, but I am lazy as fuck
//Until then I wont need it lmao
var RequestQueue []Request

func NewDormantRequest(url string) Request{
	return Request{url, nil, 0}
}

/*
Try to send a request to the preconfigured url.

returns: a request obj with a non-null return value if request is a success
*/
func TrySendRequest(url string) (Request, error){
	
	var c = &http.Client{Timeout: 10 * time.Second}
	r, err := c.Get(url)
	if err != nil {
		return Request{url, "null", -1}, err
	}
	defer r.Body.Close()
	result := &Request{url, nil, r.StatusCode}
	json.NewDecoder(r.Body).Decode(&result.ReturnValue)
	return *result, nil
}

func TryPostJoke(joke string) (Request, error){
	var url string = globals.Glob_Cache.ActiveUrl + "/api/post/{" + joke + "}"

	var c = &http.Client{Timeout: 10 * time.Second}
	r, err := c.Get(url)
	if err != nil {
		return Request{url, "null", -1}, err
	}
	defer r.Body.Close()
	result := &Request{url, nil, r.StatusCode}
	json.NewDecoder(r.Body).Decode(&result.ReturnValue)
	return *result, nil
}

/*
Get some jokes from the server and store them localy
*/
func TryGetJokes(size int) (structures.JokeList, error){
	var url string = globals.Glob_Cache.ActiveUrl + "api/fetch_jokes/{" + fmt.Sprint(size) + "}"
	var c = &http.Client{Timeout: 10 * time.Second}
	var r, err = c.Get(url)	
	if err != nil{
		return structures.JokeList{}, err
	}
	defer r.Body.Close()
	var target structures.JokeList	
	json.NewDecoder(r.Body).Decode(&target)
	return target, nil
}

/*
Retrieve the amount of jokes the database holds
*/
func GetServerJokeStock() int{
 return 0
}

/*
Tests the url to see if an applicable server is running so our requests don't get sent to russian spies and shit

NOTE: the endpoint of the url should match the testing endpoint on the server, but we'll cross that bridge once we get to it.
*/
func TestUrl(url string) bool{
	var c = &http.Client{Timeout: 4 * time.Second}
	_, err := c.Get(url)
	if err == nil {
		return true
	}
	return false
}
