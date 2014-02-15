package session

import (
    "fmt"
    //"ssdb"
    "testing"
)

func Test_test(t *testing.T) {
    ip := "127.0.0.1"
    port := 8888
    store, err := NewSSDBStore(ip, port)
    if err != nil {
        fmt.Println(err)
    }
    sid := "test_sid"
    session := Session{
        Sid:   sid,
        Store: store,
    }
    fmt.Println(session.GetSessionData())
    key := "test_key"
    value := "test_value"
    err = session.SetItem(key, value)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(session.GetSessionData())
}
