package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/beakeyz/dadjoke-gen/pkg/structures"
	. "github.com/beakeyz/dadjoke-gen/pkg/client/globals"
)

var(
	isSetup = true
)
/*
Prepare for the fun to begin :)
*/
func Prepare() (structures.JokeList, Cache, error) {

	//setup dummy vars
	var jokes structures.JokeList
	var cache Cache

	joke_file, joke_err := os.Open(Jokes_path)
	cache_file, cache_err := os.Open(Cache_path)

	//check for errors
	if joke_err != nil {
		if strings.Contains(joke_err.Error(), "no such file or directory") {

			file, _ := json.MarshalIndent(structures.JokeList{[]structures.Joke{}, 0}, "", " ")
			os.Mkdir("assets", 0755)
			_ = ioutil.WriteFile(Jokes_path, file, 0777)
			return Prepare()
		}

		fmt.Println(joke_err)
		return jokes, cache, joke_err
	} else if cache_err != nil {
		if strings.Contains(cache_err.Error(), "no such file or directory") {

			file, _ := json.MarshalIndent(Default_cache, "", " ")
			os.Mkdir("assets", 0755)
			_ = ioutil.WriteFile(Cache_path, file, 0777)
			isSetup = false
			return Prepare()
		}

		fmt.Println(cache_err)
		return jokes, cache, cache_err
	}
	//nolint:gosec

	//and parse the file
	joke_byte_arr, _ := ioutil.ReadAll(joke_file)
	cache_byte_arr, _ := ioutil.ReadAll(cache_file)

	json.Unmarshal(joke_byte_arr, &jokes)
	json.Unmarshal(cache_byte_arr, &cache)

	if !isSetup{
		isSetup = true;
		cache.HasSelf = false;
	}
	return jokes, cache, nil
}

/*
save the current local cache to a json file
*/
func SaveCacheFile(oldCache Cache) {
	file, _ := json.MarshalIndent(oldCache, "", " ")
	_ = ioutil.WriteFile(Cache_path, file, 0644)
}

/*
save the local jokes to a json file
*/
func SaveJokeFile(list structures.JokeList) {
	file, _ := json.MarshalIndent(list, "", " ")
	_ = ioutil.WriteFile(Jokes_path, file, 0422)
}

/*
brain go brrrrr, go figure
*/
func SetCache(oldCache *Cache, newCache Cache) Cache {
	oldCache = &newCache
	return *oldCache
}


