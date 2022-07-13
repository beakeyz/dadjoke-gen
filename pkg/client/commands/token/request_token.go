package token

import(
	"github.com/spf13/cobra"
)

var RequestToken = &cobra.Command{
	Use: "reqtoken",
	Run: func(cmd *cobra.Command, args []string) {
		//ask the user for email and username

		//Contact the server and ask for a token with uname and email

		//Recieve the Token from the server and store it into the cache

		//Error handling
	} ,

}
