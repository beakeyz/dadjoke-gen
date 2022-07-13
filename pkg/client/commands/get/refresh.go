package get

import (
	"fmt"

	"github.com/beakeyz/dadjoke-gen/pkg/client/requests"
	"github.com/beakeyz/dadjoke-gen/pkg/client/globals"
	"github.com/spf13/cobra"
)

var Refresh = &cobra.Command{
	Use: "refresh",
	Run: func(cmd *cobra.Command, args []string) {

		//get jokes from server
		list, err := requests.TryGetJokes(50)
		if err != nil{
			fmt.Println("Refresh failed")
			return
		}

		//Clear current jokelist
		globals.Jokes.Clear()

		//Fill jokelist with the recieved jokes
		globals.Jokes.AddJokes(list)
		fmt.Println("Jokes successfully refreshed!")
		//error handling

	},
}
