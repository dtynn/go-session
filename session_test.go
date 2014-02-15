package session

import (
    "fmt"
    "ssdb"
    "testing"
)

func Test_test(t *testing.T) {
    ip := "127.0.0.1"
    port := 8888
    client, err := ssdb.Connect(ip, port)
    if err != nil {
        fmt.Println(err)
    }
    s := SSDBStore{client, "session", 3600 * 24}
    sid := s.MakeSid()
    fmt.Println(sid)
    key := "test_key"
    data := "test_data"
    fmt.Println(s.Set(sid, key, data))
    fmt.Println(s.Get(sid, key))
    fmt.Println(s.Delete(sid))
    fmt.Println(s.Get(sid, key))
}
