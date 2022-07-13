package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/beakeyz/dadjoke-gen/pkg/structures"
)

func (self *HttpServer) Post_Jokes(rw http.ResponseWriter, req *http.Request) {

  fmt.Println("Posted!")
  
  var incommingJoke structures.Joke
  
  json.NewDecoder(req.Body).Decode(&incommingJoke)


  if incommingJoke.Summary == "" {
    rw.Write([]byte("No summary sent!"))
    return
  } else if incommingJoke.Joke == "" {
    rw.Write([]byte("No Joke sent!"))
    return
  }

  storedJokes, _err := self.GetJokesUtil()
  if _err != nil {
    fmt.Println(_err.Error())
  }

  indexToStore := storedJokes.Size

  _, err := self.db.Exec("INSERT INTO `jokes`(`Summary`, `Joke`, `AdditionDate`, `JokeIndex`) VALUES ( ?, ?, ?, ? )",
    incommingJoke.Summary,
    incommingJoke.Joke,
    incommingJoke.Date,
    indexToStore)
  if err != nil {
    fmt.Println(err.Error())
    rw.Write([]byte("Something went wrong"))
    return
  }


  rw.Write([]byte("Success!"))
}
