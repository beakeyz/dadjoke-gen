package structures

type Post struct {
  Creator User
  ImagePath string
  Caption Joke
  Collection JokeList
  CreationDate string
}
