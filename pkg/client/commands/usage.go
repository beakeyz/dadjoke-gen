package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Usage = &cobra.Command{
	Use: "usage",
	Short: "",
	Long: `Genrate dadjokes from a webserver that is connected to a database. This client is to comunicate to the webserver, which just acts as an middleman to the database. The jokes are recieved in json format and thus stored in a json file so we can easily store and access jokes localy.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`Generate dadjokes from a database that get put into a json file for fast lookuptimes
Usage:
   - Get some jokes: dadgen gen .a=[amount]
   - Get a random joke: dadgen gen .r
   - Post a joke: dadgen post .t=[token] .j=[joke]
   - Get a token: dadgen token .req
   - Remove your token: dadgen token .r`)
	},
}
