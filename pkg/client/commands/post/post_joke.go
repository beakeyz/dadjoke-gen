package post

import(
	"github.com/spf13/cobra"
)

var PostJoke = &cobra.Command{
	Use: "post",
	Run: func(cmd *cobra.Command, args []string) {
		//look for the -j flag in order to parse the joke

		//parse the joke

		//make a request to store the joke

		//recieve status code from the server to see if the joke was posted successfully or if it's pending 

		//Error handling
	},
}
