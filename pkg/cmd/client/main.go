package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beakeyz/dadjoke-gen/pkg/client/commands"
	"github.com/beakeyz/dadjoke-gen/pkg/client/commands/get"
	"github.com/beakeyz/dadjoke-gen/pkg/client/commands/token"
	"github.com/beakeyz/dadjoke-gen/pkg/client/commands/post"
	"github.com/beakeyz/dadjoke-gen/pkg/client/globals"
	"github.com/beakeyz/dadjoke-gen/pkg/client/utils"
	"github.com/beakeyz/dadjoke-gen/pkg/client/requests"
)

//What do you call someone with no body and no nose? Nobody knows.

func main() {
	fmt.Println("Module successfully started!")
	/* setup */

	var err error
	//Retrieve the current data from the cache, if it is valid
	globals.Jokes, globals.Glob_Cache, err = utils.Prepare()
	if (err != nil){
		fmt.Println("Something went wrong while obtaining locales")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if !globals.Glob_Cache.HasSelf{
		fmt.Println("Looks like ur install aint correctly setup...")
		fmt.Println("What is the url that you plan to use?")
		fmt.Println("input: ")
		var url string
		for {
			
			fmt.Scanf("%s", &url)

			if (url != "" && requests.TestUrl(url)){
				globals.Glob_Cache.ActiveUrl = url
				globals.Glob_Cache.HasSelf = true;
			}else{
				fmt.Println("That url is not gay enough")
				utils.SaveCacheFile(globals.Glob_Cache)
				os.Exit(-1)
			}	

			if (globals.Glob_Cache.HasSelf){
				break
			}
			time.Sleep(1)
		}
	}

	/* prepare the local cache */
	globals.Glob_Cache.UpdateCacheDate()

	/* do the module functionality here */

	PrepareCmd()

	/* save the local cache to json format */

	globals.Jokes.RunJokeCheckAndSort()
	globals.Glob_Cache.CorrectCache()
	//globals.Glob_Cache.UpdateLocalCache()
	utils.SaveJokeFile(globals.Jokes)
	utils.SaveCacheFile(globals.Glob_Cache)
}

func PrepareCmd(){

	commands.Rootcmd.AddCommand(get.Get)
	commands.Rootcmd.AddCommand(get.Refresh)
	commands.Rootcmd.AddCommand(commands.Usage)
	commands.Rootcmd.AddCommand(post.PostJoke)
	commands.Rootcmd.AddCommand(token.RequestToken)

	if err := commands.Rootcmd.Execute(); err != nil{
		fmt.Println("oof, something went wrong lmao")
		os.Exit(1)
	}
}
