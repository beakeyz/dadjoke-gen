package get

import (
	"github.com/beakeyz/dadjoke-gen/pkg/client/globals"
	"github.com/spf13/cobra"
	"fmt"
	"math/rand"
)

func Random(cmd *cobra.Command, args []string){
	if globals.Jokes.Size == 0{
		fmt.Println("You don't have any jokes: get them using the `gen amount` command")
		return
	}
	end := globals.Jokes.Size - 1
	random := rand.Intn(end)	
	joke := globals.Jokes.GetJoke(0 + random)
	globals.Glob_Cache.PreviousJokeIndex = joke.Index
	fmt.Println(joke.Joke)
}
