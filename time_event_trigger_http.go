package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func (tet *TimeEventTrigger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		tet.servePost(w, r)
		return
	default:
		w.Write([]byte(REP_ROUTER_NOT_FOUND))
		return
	}
}

func (tet *TimeEventTrigger) servePost(w http.ResponseWriter, r *http.Request) {
	var reqObj struct {
		Key   string `json:"key"`
		Data  string `json:"data"`
		Delay int64  `json:"delay"`
	}
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(REP_BAD_REQ))
		return
	}
	err = json.Unmarshal(buf, &reqObj)
	if err != nil {
		w.Write([]byte(REP_BAD_REQ))
		return
	}
	if reqObj.Key == "" || reqObj.Delay < 0 {
		w.Write([]byte(REP_NO_KEY))
		return
	}
	err = tet.addTimeEvent(&TimeEvent{
		key:   reqObj.Key,
		data:  []byte(reqObj.Data),
		delay: time.Second * time.Duration(reqObj.Delay),
	})
	if err != nil {
		w.Write([]byte("add event error"))
		return
	}
	w.Write([]byte(REP_OK))
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
