package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/tj/go-debug"
)

var __deRegMgr__ = Debug("timer:registerMgr")

type RegisterMgr struct {
	clients           []RegisterClient
	delayedTimeEvents []*DelayedTimeEvent
}

func NewRegisterMgr() *RegisterMgr {
	return &RegisterMgr{
		clients:           make([]RegisterClient, 0),
		delayedTimeEvents: make([]*DelayedTimeEvent, 0),
	}
}

func (rgm *RegisterMgr) HandleTimeEvent(te *TimeEvent) error {
	__deRegMgr__("got time event to handle:%s", te)
	key := te.key
	for _, c := range rgm.clients {
		if c.Match(key) {
			err := c.Deliver(te)
			if err != nil {
				rgm.delay(te, c)
				continue
			}
		}
	}
	return nil
}

func (rgm *RegisterMgr) delay(te *TimeEvent, cli RegisterClient) {
	__deRegMgr__("clients not found or unavaliable for te:%s with name %s. delay")
	rgm.delayedTimeEvents = append(rgm.delayedTimeEvents, &DelayedTimeEvent{
		client: cli,
		te:     te,
	})
}

func (rgm *RegisterMgr) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		rgm.Post(w, r)
		return
	default:
		w.Write([]byte(REP_ROUTER_NOT_FOUND))
	}
}

//add http client
//err if name exists
func (rgm *RegisterMgr) Post(w http.ResponseWriter, r *http.Request) {
	//key name url isPattern
	var reqObj struct {
		Name      string `json:"name"`
		Url       string `json:"url"`
		Key       string `json:"key"`
		isPattern bool   `json:"isPattern"`
	}
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(REP_BAD_REQ))
		return
	}
	err = json.Unmarshal(buf, &reqObj)
	if err != nil {
		w.Write([]byte("err parse json:" + err.Error()))
		return
	}
	__deRegCli__("serve http for register mgr:%v", reqObj)
	if reqObj.Key == "" {
		w.Write([]byte(("no key")))
		return
	}
	if reqObj.Name == "" {
		w.Write([]byte("no name"))
		return
	}
	//TODO regexp url
	if reqObj.Key == "" {
		w.Write([]byte("no url or illegal"))
		return
	}
	var m RegisterClientMatcher
	if reqObj.isPattern == false {
		m = NewFixedKeyMatcher(reqObj.Key)
	} else {
		var errCpl error
		m, errCpl = NewRegexpKeyMatcher(reqObj.Key)
		if errCpl != nil {
			w.Write([]byte(fmt.Sprintf("compile regexp:%s, fail, err:%s", reqObj.Key, errCpl)))
			return
		}
	}
	c := NewRegisterClientHTTP(m, reqObj.Name, reqObj.Url)
	rgm.clients = append(rgm.clients, c)
	__deRegMgr__("add http client, key:%s, name:%s, url:%s",
		reqObj.Key, reqObj.Name, reqObj.Url)
	w.Write([]byte(REP_OK))
}

type DelayedTimeEvent struct {
	client RegisterClient
	te     *TimeEvent
}
