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
	Sessions    []Session
	SessionPath string
}

type Session struct {
	LinkedUser   User
	CreationDate string
	MaxAge       int
	SessionId    uuid.UUID
	FileName     os.File
	IsNull       bool
}

func EmptySession() *Session {
	return &Session{IsNull: true}
}

func (self *Session) SetSession(newSass *Session) {
	self = newSass
}

func CreateSessionTemplate(user *User, sassId uuid.UUID) *Session {

	return &Session{
		LinkedUser:   *user,
		CreationDate: time.Now().Format(time.RFC3339),
		MaxAge:       int(time.Hour * time.Duration(24)),
		SessionId:    sassId,
		FileName:     os.File{},
	}
}

func createSession(sass *Session) (*Session, error) {
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
		Sessions:    []Session{},
		SessionPath: sessionPath,
	}

	if refreshErr := RefreshSessions(manager); refreshErr != nil {
		fmt.Println(refreshErr.Error())
		return &SessionManager{}, refreshErr
	}

	return manager, nil
}

func RefreshSessions(mngr *SessionManager) error {
	var items, _ = ioutil.ReadDir(sessionPath)
	for _, item := range items {
		//fmt.Println("Session: " + item.Name())
		if !item.IsDir() && strings.Contains(item.Name(), ".json") {
			var raw, _ = os.Open(strings.Join([]string{sessionPath, item.Name()}, ""))
			var fileBytes, readErr = ioutil.ReadAll(raw)
			if readErr != nil {
				return readErr
			}
			var dummySession Session = Session{}
			json.Unmarshal(fileBytes, &dummySession)

			// now we can perform checks on the session

			sessTime, parseErr := time.Parse(time.RFC3339, dummySession.CreationDate)
			if parseErr != nil {
				// cry once again
				fmt.Println(parseErr.Error())
			}

			if time.Since(sessTime) > time.Duration(dummySession.MaxAge) {
				//fmt.Println("Session expired!")
				// remove session
				if removeErr := mngr.RemoveSession(&dummySession); removeErr != nil {
					return removeErr
				}
			} else {
				mngr.Sessions = append(mngr.Sessions, dummySession)
			}
		}
	}
	return nil
}

// TODO perhaps have an RefreshSessions function that syncs the local sessions in memory with the sessions on disk?
func (self *SessionManager) AddSession(session *Session) error {
	if refreshErr := RefreshSessions(self); refreshErr != nil {
		fmt.Println(refreshErr.Error())
		return refreshErr
	}

	sass, sassErr := createSession(session)
	if sassErr != nil {
		fmt.Println(sassErr.Error())
		return sassErr
	}
	self.Sessions = append(self.Sessions, *sass)
	return nil
}

// Check for expired sessions and delete them
func (self *SessionManager) ClearSessions() {

}

func (self *SessionManager) GetSessionFromUser(user *User) (*Session, error) {
	if refreshErr := RefreshSessions(self); refreshErr != nil {
		fmt.Println(refreshErr.Error())
		return EmptySession(), refreshErr
	}

	for _, sass := range self.Sessions {
		// NOTE: dereference the user param, bcuz we need to check if the objects are the same, not if their addresses match =D
		if sass.LinkedUser == *user {
			return &sass, nil
		}
	}
	return EmptySession(), nil
}

func (self *SessionManager) GetSession(Uuid uuid.UUID) (*Session, error) {

	if refreshErr := RefreshSessions(self); refreshErr != nil {
		fmt.Println(refreshErr.Error())
		return EmptySession(), refreshErr
	}

	for _, sass := range self.Sessions {
		if sass.SessionId == Uuid {
			return &sass, nil
		}
	}
	return EmptySession(), nil
}

func (self *SessionManager) RemoveSession(sass *Session) error {
  fmt.Println("pulled a quick one")
	if self.ContainsSession(sass) {
		// TODO: do funnie and yeet the session
		var sessionIndex int = 0
		for index, thing := range self.Sessions {
			if thing == *sass {
        sessionIndex = index
        break
			}
		}
    self.Sessions = append(self.Sessions[:sessionIndex], self.Sessions[sessionIndex + 1:]...)
	}

	removeErr := os.Remove(sessionPath + sass.SessionId.String() + ".json")
	if removeErr != nil {
    fmt.Println("fucked up while removing a session.json: " + removeErr.Error())
		return removeErr
	}

	return nil
}

func (self *SessionManager) ContainsSession(sess *Session) bool {
	for _, b := range self.Sessions {
		if b == *sess {
			return true
		}
	}
	return false
}
