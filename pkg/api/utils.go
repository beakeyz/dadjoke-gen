package api

import (
	"fmt"

	"github.com/beakeyz/dadjoke-gen/pkg/structures"
)

func (self *HttpServer) GetJokesUtil() (structures.JokeList, error) {
  var jokes structures.JokeList

	rows, err := self.db.Query("SELECT * FROM jokes WHERE 1")
	if err != nil {
    fmt.Println(err.Error())
    return jokes, err
	}
	defer rows.Close()

	for rows.Next() {
		jokes.Size++
		if jokes.Size > 10 {
			break
		}
		var tmp structures.Joke
		if err := rows.Scan(&tmp.Summary, &tmp.Joke, &tmp.Date, &tmp.Index); err != nil {
      fmt.Println(err.Error())
      return jokes, err
		}

		jokes.List = append(jokes.List, tmp)
	}

  return jokes, nil
}
