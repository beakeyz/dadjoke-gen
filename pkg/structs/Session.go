package structs

import (
	"container/list"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
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
		IsNull:       false,
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

/* SESSION MANAGER (mostly lol) */

// TODO: move these into SessionManager
var (
	sessionRequestList  *list.List = list.New()
	sessionResponseList *list.List = list.New()
)

type SESSION_REQUEST_TYPE int8

const (
  session_request_remove SESSION_REQUEST_TYPE = 0
  session_request_add SESSION_REQUEST_TYPE = 1
  session_request_replace SESSION_REQUEST_TYPE = 2
  session_request_get_usr SESSION_REQUEST_TYPE = 3
  session_request_get_uuid SESSION_REQUEST_TYPE = 4
)

type SessionEditRequest struct {
	Name         string
	Request_type SESSION_REQUEST_TYPE 

	Sass Session
	Usr  User
	Uuid uuid.UUID

	Id int64
}

func prepSessionRequest(session Session, req_type SESSION_REQUEST_TYPE) SessionEditRequest {
	req := SessionEditRequest{}
	one, two := rand.Int(rand.Reader, big.NewInt(9999998))
	if two != nil {
		one = big.NewInt((int64(sessionRequestList.Len() + 1)))
	}
	fmt.Printf("set id %d \n", one.Int64())
	req.Id = one.Int64()
	req.Name = "New Request"
	req.Sass = session
	req.Usr = session.LinkedUser
	req.Uuid = session.SessionId

	if !(req_type >= 0 && req_type <= 4) {
		req.Request_type = session_request_get_usr
	}
	req.Request_type = req_type

	return req
}

type SessionEditResponse struct {
	Id            int64
	ReturnSession Session
	err           error
}

// Blocking requests that can be sent
func SessionRequestRemove(session Session) {
	fmt.Println("Tried to remove session")
	req := prepSessionRequest(session, session_request_remove)
	sessionRequestList.PushBack(req)
	waitForResponse(req)
}

func SessionRequestAdd(session Session) {
	fmt.Println("Tried to add session")
	req := prepSessionRequest(session, session_request_add)
	sessionRequestList.PushBack(req)

  response := waitForResponse(req)

  if response.err != nil {
    fmt.Println(response.err.Error())
  }
}

func SessionRequestReplace(session Session) {
	fmt.Println("Tried to replace session")
	sessionRequestList.PushFront(prepSessionRequest(session, session_request_replace))
}

func SessionRequestGetFromUser(sessionUser User) (Session, error) {
	// TODO
	s := *EmptySession()
	s.LinkedUser = sessionUser
	req := prepSessionRequest(s, 3)
	sessionRequestList.PushBack(req)
	var response SessionEditResponse = waitForResponse(req)

	if response.Id != -1 {
		return response.ReturnSession, response.err
	}
	return *EmptySession(), fmt.Errorf("Failed to get session by user")
}

func SessionRequestGet(session Session) (Session, error) {
	fmt.Println("Tried to get session")
	req := prepSessionRequest(session, session_request_get_uuid)
	sessionRequestList.PushBack(req)
	var res SessionEditResponse = waitForResponse(req)

	if res.Id != -1 {
		return res.ReturnSession, res.err
	}

	return *EmptySession(), fmt.Errorf("Invalid responseID recieved from SessionRequestGet!")
}

func waitForResponse(req SessionEditRequest) SessionEditResponse {
	var response SessionEditResponse = SessionEditResponse{}
	fmt.Printf("waiting... %d \n", req.Id)
	for true {
		for i := sessionResponseList.Front(); i != nil; i = i.Next() {
			fmt.Println("dummy find")
			//if i != nil {
			_res := i.Value.(SessionEditResponse)
			if _res.Id == req.Id {
				response = _res
				fmt.Println("Got a response =D")
				sessionResponseList.Remove(i)
				return response
			}
			//}
		}
		time.Sleep(2 * time.Millisecond)
	}
	// cant happen lol
	return SessionEditResponse{-1, *EmptySession(), fmt.Errorf("Returned early from waitForResponse")}
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

	if refreshErr := manager.RefreshSessions(manager); refreshErr != nil {
		fmt.Println(refreshErr.Error())
		return &SessionManager{}, refreshErr
	}

	return manager, nil
}

// Should only be called once at startup, idealy
func (self *SessionManager) RefreshSessions(mngr *SessionManager) error {
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

func (self *SessionManager) Run(ctx context.Context) error {

	for true {
		if sessionRequestList.Len() != 0 {
			// we have a session =D
      // TODO: add mutexes for extra safety (just so we know for sure that no one else tries to yoink our requests from somewhere else)

			var front list.Element = *sessionRequestList.Front()
			if front.Value == nil {
				fmt.Println("Had to remove null entry in queue")
				sessionRequestList.Remove(&front)
				continue
			}
			var request SessionEditRequest = front.Value.(SessionEditRequest)
			var _session_to_return Session = *EmptySession()
			var _err_to_return error = nil

			// do crap
			switch request.Request_type {
			// remove
			case session_request_remove:
				self.RemoveSession(&request.Sass)
				break
			// add
			case session_request_add:
				self.AddSession(&request.Sass)
				break
			// replace
			case session_request_replace:
				sass, get_err := self.GetSessionFromUser(&request.Sass.LinkedUser)
				if get_err != nil {
					// yikes
				}
				if self.RemoveSession(sass) == nil {
					self.AddSession(&request.Sass)
				}
				break
			// TODO
			// get (user)
			case session_request_get_usr:
				thing, err := self.GetSessionFromUser(&request.Sass.LinkedUser)
				_session_to_return = *thing
				_err_to_return = err

				break
			// get (uuid)
			case session_request_get_uuid:
				fmt.Println("Getting the session...")
				thing, err := self.GetSession(request.Sass.SessionId)
				_session_to_return = *thing
				_err_to_return = err

				break
			}

			fmt.Printf("Pushed back a Response =D %d \n", _session_to_return.SessionId.ID())
			// "return" a response
			sessionResponseList.PushFront(SessionEditResponse{request.Id, _session_to_return, _err_to_return})
			sessionRequestList.Remove(&front)
		}
		// sleep
		time.Sleep(1 * time.Millisecond)
	}

	// no >=(
	fmt.Println("recursed into Run due to weird exit")
	self.Run(ctx)
	return nil
}

// TODO perhaps have an RefreshSessions function that syncs the local sessions in memory with the sessions on disk?
func (self *SessionManager) AddSession(session *Session) error {
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
	if user == nil {
		return EmptySession(), fmt.Errorf("Passed a Nil user!")
	}

	for _, sass := range self.Sessions {
		// NOTE: dereference the user param, bcuz we need to check if the objects are the same, not if their addresses match =D
		if sass.LinkedUser.Token == user.Token {
			return &sass, nil
		}
	}
	return EmptySession(), nil
}

func (self *SessionManager) GetSession(Uuid uuid.UUID) (*Session, error) {

	for _, sass := range self.Sessions {
		if sass.SessionId == Uuid {
			return &sass, nil
		}
	}
	return EmptySession(), nil
}

func (self *SessionManager) RemoveSession(sass *Session) error {

	if self.ContainsSession(sass) {
		// TODO: do funnie and yeet the session
		var sessionIndex int = 0
		for index, thing := range self.Sessions {
			if thing == *sass {
				sessionIndex = index
				break
			}
		}
		self.Sessions = append(self.Sessions[:sessionIndex], self.Sessions[sessionIndex+1:]...)
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
