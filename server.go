package main

import (
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	. "github.com/tj/go-debug"
)

var __deMain__ = Debug("timer:main")

func main() {
	chTimeEventOut := make(chan *TimeEvent)
	chTimeEventCmd := make(chan *TimeEventCmd, 10)
	trigger := NewTimeEventTrigger(chTimeEventOut, chTimeEventCmd)
	go trigger.Begin()
	debugTecm := NewDebugTECM()
	regMgrTecm := NewRegisterMgr()
	go broadcastEvent(chTimeEventOut, debugTecm, regMgrTecm)
	//http register client
	//beego.Post("/time_event/delay", func(ctx *context.Context) {
	//	var reqObj struct {
	//		Key   string `json:"key"`
	//		Data  string `json:"data"`
	//		Delay int64  `json:"delay"`
	//	}
	//	buf, err := ioutil.ReadAll(r.Body)
	//	if err != nil {
	//		w.Write([]byte(REP_BAD_REQ))
	//		return
	//	}
	//	err = json.Unmarshal(buf, &reqObj)
	//	if err != nil {
	//		w.Write([]byte(REP_BAD_REQ))
	//		return
	//	}
	//	if reqObj.Key == "" || reqObj.Delay < 0 {
	//		w.Write([]byte(REP_NO_KEY))
	//		return
	//	}
	//	err = tet.AddTimeEvent(&TimeEvent{
	//		key:   reqObj.Key,
	//		data:  []byte(reqObj.Data),
	//		delay: time.Second * time.Duration(reqObj.Delay),
	//	})
	//	if err != nil {
	//		w.Write([]byte("add event error"))
	//		return
	//	}
	//	w.Write([]byte(REP_OK))
	//})
	beego.Post("/time_event/interval", func(ctx *context.Context) {
		var reqObj struct {
			Key       string `json:"key"`
			Data      string `json:"data"`
			Schedule  string `json:"schedule"`
			Listeners []struct {
				Url  string `json:"url"`
				Name string `json:"name"`
			} `json:"listeners"`
		}
		err := json.Unmarshal(ctx.Input.RequestBody, &reqObj)
		if err != nil {
			ctx.Output.Body([]byte(REP_BAD_REQ))
			return
		}
		m := NewTimeEventMeta(reqObj.Key, []byte(reqObj.Data))
		interGen, err := newIntervalTimeEventGenerator(reqObj.Schedule, m)
		trigger.AddTimeEventGenerator(interGen)
		//add listener
		matcher := NewFixedKeyMatcher(reqObj.Key)
		for _, l := range reqObj.Listeners {
			regMgrTecm.AddClients(NewRegisterClientHTTP(matcher, l.Name, l.Url))
		}
	})
	//http.Handle("/handler/http", regMgrTecm)
	http.ListenAndServe(":6200", nil)
}

func broadcastEvent(chTE <-chan *TimeEvent, heads ...TimeEventHandler) {
	for {
		__deMain__("wait for te")
		te := <-chTE
		for _, head := range heads {
			err := head.HandleTimeEvent(te)
			if err != nil {
				__deMain__("handle error:%s, quit", err)
				break
			}
		}
	}
}

const pathPrefix = "time_event"
const KEY_NAME = "key"
const DATA_NAME = "data"
const DELAY_NAME = "delay"

var (
	REP_ROUTER_NOT_FOUND = "404 router not found\n"
	REP_OK               = "OK\n"
	REP_DELAY_WRONG      = "param delay must be posive(in seconds)\n"
	REP_NO_KEY           = "no key\n"
	REP_BAD_REQ          = "request body or query error\n"
)
