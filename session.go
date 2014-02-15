package session

import (
    "crypto/md5"
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

//session store interface
type SessionStore interface {
    Get(sid string, key string) (interface{}, error)
    Set(sid string, key string, data string) (interface{}, error)
    Delete(sid string) (interface{}, error)
}

type SSDBStore struct {
    client *ssdb.Client
    prefix string
    expire int32
}

func (s *SSDBStore) MakeSid() string {
    return fmt.Sprintf("%s:%s", s.prefix, GenerateUUID())
}

func (s *SSDBStore) Get(sid string, key string) (interface{}, error) {
    resp, err := s.client.Do("hget", sid, key)
    if err != nil {
        return nil, err
    }
    if len(resp) > 2 {
        err := errors.New("bad response")
        return nil, err
    }
    if resp[0] != "ok" {
        err := errors.New(resp[0])
        return nil, err
    }
    return resp[1], nil
}

func (s *SSDBStore) Set(sid string, key string, data string) (interface{}, error) {
    resp, err := s.client.Do("hset", sid, key, data)
    if err != nil {
        return nil, err
    }
    if len(resp) > 2 {
        err := errors.New("bad response")
        return nil, err
    }
    if resp[0] != "ok" {
        err := errors.New(resp[0])
        return nil, err
    }
    return resp[1], nil
}

func (s *SSDBStore) Delete(sid string) (interface{}, error) {
    resp, err := s.client.Do("hclear", sid)
    if err != nil {
        return nil, err
    }
    if len(resp) > 2 {
        err := errors.New("bad response")
        return nil, err
    }
    if resp[0] != "ok" {
        err := errors.New(resp[0])
        return nil, err
    }
    return resp[1], nil
}

//Session
type Session struct {
    sid   string
    store *SessionStore
}

func (s *Session) clear() {
    s.store.Delete(s.sid)
}
