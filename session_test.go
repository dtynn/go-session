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
		return
	}
	sid := "test_sid"
	session := NewSession(sid, store)
	fmt.Println(session.Sid)
	fmt.Println(session.IsNew)
	err = session.GetSession()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(session.Values)
	session.Save()
	fmt.Println(session.Values)
	key := "test_key_" + GenerateUUID()
	value := "test_value_" + GenerateUUID()
	session.SetItem(key, value)
	fmt.Println(session.Values)
	fmt.Println(session.GetItem(key))
	fmt.Println(session.Contains(key))
	fmt.Println(session.ClearSession())
	fmt.Println(session.Values)
	fmt.Println(session.GetItem(key))
	fmt.Println(session.Contains(key))
	//fmt.Println(session.GetSessionData())
	//key := "test_key"
	//value := "test_value"
	//err = session.SetItem(key, value)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//resMap, _ := session.GetSessionData()
	//fmt.Println(resMap)
	//fmt.Println(resMap["test_key"])
	//fmt.Println(session.Clear())
	//fmt.Println(session.GetSessionData())
	return
}
