package session

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"ssdb"
	"strconv"
	"time"
)

func GenerateUUID() string {
	nano := time.Now().UnixNano()
	r := rand.New(rand.NewSource(nano))
	num := r.Int63()
	mixed := GenerateMD5(strconv.FormatInt(nano, 10)) + GenerateMD5(strconv.FormatInt(num, 10))
	return GenerateMD5(mixed)
}

func GenerateMD5(text string) string {
	hashed := md5.New()
	io.WriteString(hashed, text)
	return fmt.Sprintf("%x", hashed.Sum(nil))
}

func MakeSid(prefix string) string {
	return fmt.Sprintf("%s:%s", prefix, GenerateUUID())
}

//session store interface
type SessionStore interface {
	Get(sid string, key string) (string, error)
	Set(sid string, key string, data string) error
	Del(sid string, key string) error
	Clear(sid string) error
}

type SSDBStore struct {
	client *ssdb.Client
}

func (s *SSDBStore) Get(sid string, key string) (res string, err error) {
	resp, err := s.client.Do("hget", sid, key)
	if err != nil {
		return res, err
	}
	if len(resp) > 2 {
		err := errors.New("bad response")
		return res, err
	}
	if resp[0] == "not_found" {
		return res, nil
	} else if resp[0] != "ok" {
		err := errors.New(resp[0])
		return res, err
	}
	return resp[1], nil
}

func (s *SSDBStore) Set(sid string, key string, data string) error {
	resp, err := s.client.Do("hset", sid, key, data)
	if err != nil {
		return err
	}
	if len(resp) > 2 {
		err := errors.New("bad response")
		return err
	}
	if resp[0] != "ok" {
		err := errors.New(resp[0])
		return err
	}
	return nil
}

func (s *SSDBStore) Del(sid string, key string) error {
	resp, err := s.client.Do("hdel", sid, key)
	if err != nil {
		return err
	}
	if len(resp) > 2 {
		err := errors.New("bad response")
		return err
	}
	if resp[0] != "ok" {
		err := errors.New(resp[0])
		return err
	}
	return nil
}

func (s *SSDBStore) Clear(sid string) error {
	resp, err := s.client.Do("hclear", sid)
	if err != nil {
		return err
	}
	if len(resp) > 2 {
		err := errors.New("bad response")
		return err
	}
	if resp[0] != "ok" {
		err := errors.New(resp[0])
		return err
	}
	return nil
}

func NewSSDBStore(ip string, port int) (*SSDBStore, error) {
	client, err := ssdb.Connect(ip, port)
	if err != nil {
		return nil, err
	}
	return &SSDBStore{client}, nil
}

//Session
type Session struct {
	Sid    string
	Store  SessionStore
	Values map[string]interface{}
	IsNew  bool
	//Expire int32
}

func (s *Session) GetSessionData() error {
	data, err := s.Store.Get(s.Sid, "data")
	if err != nil {
		return err
	}
	if data != "" {
		err = json.Unmarshal([]byte(data), &s.Values)
		if err != nil {
			return err
		}
	} else {
		s.Values = nil
	}
	s.IsNew = true
	return nil
}

func (s *Session) GetSession() error {
	if s.IsNew == true {
		return nil
	}
	err := s.GetSessionData()
	return err
}

func (s *Session) SetItem(key string, value interface{}) error {
	s.Dirty()
	if s.Values == nil {
		s.Values = make(map[string]interface{})
	}
	s.Values[key] = value
	err := s.Save()
	return err
}

func (s *Session) GetItem(key string) (interface{}, error) {
	err := s.GetSession()
	if err != nil {
		return nil, err
	}
	for k, v := range s.Values {
		if k == key {
			return v, nil
		}
	}
	return nil, nil
}

func (s *Session) Contains(key string) (bool, error) {
	err := s.GetSession()
	if err != nil {
		return false, err
	}
	for k, _ := range s.Values {
		if k == key {
			return true, nil
		}
	}
	return false, nil
}

func (s *Session) ClearSession() error {
	s.Dirty()
	err := s.Store.Clear(s.Sid)
	if err != nil {
		return err
	}
	err = s.GetSessionData()
	return err
}

func (s *Session) Dirty() {
	s.IsNew = false
}

func (s *Session) Save() error {
	b, err := json.Marshal(s.Values)
	if err != nil {
		return err
	}
	err = s.Store.Set(s.Sid, "data", string(b))
	if err != nil {
		return err
	}
	s.GetSessionData()
	return nil
}

func NewSession(sid string, store SessionStore) *Session {
	return &Session{
		Sid:   sid,
		Store: store,
	}
}
