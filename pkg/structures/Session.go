package structures

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
  sessionPath = "sessions/"
)

type SessionManager struct {
  Sessions []Session
  SessionPath string
}

type Session struct {
  LinkedUser User
  ExperationDate string
  SessionId uuid.UUID
  FileName os.File
  IsNull bool
}

func EmptySession () *Session {
  return &Session{IsNull: true}
}

func createSession(user *User) *Session {
  return &Session{

  }
}

func createAnonymousSession(anonUser *User) (*Session, error) {

  var sass *Session = &Session {
    LinkedUser: *anonUser,
    ExperationDate: time.Now().Add(time.Hour * 24).Format("2022-07-12"),
    SessionId: anonUser.Token,
    FileName: os.File{},
  }

  bytes, err := json.Marshal(sass)
  if err != nil {
    return &Session{}, err
  }

  var fileName string = sessionPath + strings.Join([]string{sass.SessionId.String(), "json"}, ".")
  var sassFileError error = os.WriteFile(fileName, bytes, 0777)
  if sassFileError != nil {
    return &Session{}, sassFileError
  }

  var sassFile, openError = os.Open(fileName)
  if openError != nil {
    return &Session{}, openError
  }

  sass.FileName = *sassFile

  return sass, nil

}

func CreateSassManager() (*SessionManager, error) {
  var err error = os.MkdirAll(sessionPath, 0777)
  if err != nil {
    return &SessionManager{}, err 
  }

  var manager *SessionManager = &SessionManager{
    Sessions: []Session{},
    SessionPath: sessionPath,
  }

  if refreshErr := RefreshSessions(manager); refreshErr != nil {
    fmt.Println(refreshErr.Error())
    return &SessionManager{}, refreshErr
  }

  return manager, nil
}

func RefreshSessions (mngr *SessionManager) error {
  var items, _ = ioutil.ReadDir(sessionPath) 
  for _, item := range items {
    fmt.Println("Hi: " + item.Name())
    if !item.IsDir() && strings.Contains(item.Name(), ".json") {
      var raw, _ = os.Open(strings.Join([]string{sessionPath, item.Name()}, ""))
      var fileBytes, readErr = ioutil.ReadAll(raw)
      if readErr != nil {
        return readErr
      }
      var dummySession Session = Session{}
      json.Unmarshal(fileBytes, &dummySession)
      mngr.Sessions = append(mngr.Sessions, dummySession)
    }
  } 
  return nil
}

// TODO perhaps have an RefreshSessions function that syncs the local sessions in memory with the sessions on disk?
func (self *SessionManager) AddSession(user *User) error {
  if refreshErr := RefreshSessions(self); refreshErr != nil {
    fmt.Println(refreshErr.Error())
    return refreshErr
  }

  if user.IsAnonymous {
    sass, sassErr := createAnonymousSession(user)
    if sassErr != nil {
      fmt.Println(sassErr.Error())
      return sassErr 
    } 
    self.Sessions = append(self.Sessions, *sass)
    return nil
  }
  return nil
}

// Check for expired sessions and delete them
func (self *SessionManager) ClearSessions() {

}

func (self *SessionManager) GetSession(user *User) (*Session, error) {
  if refreshErr := RefreshSessions(self); refreshErr != nil {
    fmt.Println(refreshErr.Error())
    return EmptySession(), refreshErr
  }

  for _, sass := range self.Sessions {
    if sass.LinkedUser.Token == user.Token {
      return &sass, nil
    } 
  }
  return EmptySession(), nil
}
