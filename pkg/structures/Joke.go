package structures

import (
	"strings"
	"time"
)

type Joke struct {
  Summary string
	Joke  string
  Date  string
	Index int
}

type JokeList struct {
	List           []Joke
	Size           int
}

func Make_jokelist(list []Joke) JokeList {
	var jokeList JokeList
	jokeList.List = list
	jokeList.Size = len(list)
	return jokeList
}

func (self *JokeList) Clear(){
	self.List = nil 
	self.Size = 0
}

func (list JokeList) GetJoke(index int) Joke {
	return list.List[index]
}

func (list JokeList) GetJokeWithInx(index int) (int, Joke) {
	i := 0
	for _, joke := range list.List {
		if joke.Index == index {
			return i, joke
		}
		i++
	}
	return -1, Joke{}
}

func (list JokeList) GetJokeStr(str string) Joke {
	i := 0
	for range list.List {
		if strings.Contains(list.List[i].GetJoke(), str) {
			return list.GetJoke(i)
		}
		i++
	}
	return Joke{}
}

func (self *JokeList) RunJokeCheckAndSort() JokeList {
	//TODO: check the jokes and sort based on index they have as property vs. what it should be in the array.

	//var indices []int
	i := 0
	for _, joke := range self.List {

		self.List[i] = Joke{joke.Summary, joke.GetJoke(), joke.Date, i + 1}

		//fmt.Println(self.List[i].Joke)
		//fmt.Println(joke.Index)
		i++
	}

	self.Size = i

	return *self
}

func (list *JokeList) AddJoke(summary string, joke string) {
	list.List = append(list.List, Joke{summary, joke, time.Now().Format("yyyy-mm-dd"), 0})
	list.RunJokeCheckAndSort()
}

func (self *JokeList) AddJokes(list JokeList){
	for _, Joke := range list.List{
		self.AddJoke(Joke.Summary, Joke.Joke)	
	}
}

/*
remove a joke from the json storage, using the Index attribute that every joke has
*/
func (list *JokeList) RemoveJoke(index int) {
	i, _ := list.GetJokeWithInx(index)
	list.RunJokeCheckAndSort()
	list.List = append(list.List[:i], list.List[i+1:]...)
	list.RunJokeCheckAndSort()
}

func (j Joke) GetJoke() string {
	return j.Joke
}
