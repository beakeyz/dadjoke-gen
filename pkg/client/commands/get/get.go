package get

import (
	"github.com/spf13/cobra"
)

var Get = &cobra.Command{
	Use: "gen",
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			if arg == ".r" || arg == "random"{
				Random(cmd, args)
				return
			}else if arg == ".a"{
				
			}
		}
	},
}
