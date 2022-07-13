package commands

import(
	"github.com/spf13/cobra"
)

var Rootcmd = &cobra.Command{
	Use: "dadgen",
	Short: "Generate dadjokes from a database that get put into a json file for fast lookuptimes",
	Long: `Genrate dadjokes from a webserver that is connected to a database. This client is to comunicate to the webserver, which just acts as an middleman to the database. The jokes are recieved in json format and thus stored in a json file so we can easily store and access jokes localy.`,
}
