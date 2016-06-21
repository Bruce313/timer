package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/tj/go-debug"
)

var (
	__deTrigger__ = debug.Debug("timer:trigger")
)

//TimeEventTrigger is a wrapper which triggers timeevent
type TimeEventTrigger struct {
	chTimeEventOut chan<- *TimeEvent
	chTimeEventCmd <-chan *TimeEventCmd
	//TODO heap it
	tes     []*TimeEvent
	timer   *time.Timer
	nearest *TimeEvent
	beego.Controller
}

//NewTimeEventTrigger create TimeEventTrigger
func NewTimeEventTrigger(ch chan<- *TimeEvent, chType <-chan *TimeEventCmd) *TimeEventTrigger {
	timer := time.NewTimer(1 * time.Hour)
	timer.Stop()
	return &TimeEventTrigger{
		chTimeEventOut: ch,
		chTimeEventCmd: chType,
		timer:          timer,
		tes:            make([]*TimeEvent, 0),
	}
}

// //Add append timeevent to a key
// func (tet *TimeEventTrigger)Add(te *TimeEvent) {
//     tet.tes = append(tet, te)
// }

//Begin start wait timeevent and trigger it
func (tet *TimeEventTrigger) Begin() {
	for {
		select {
		case cmd := <-tet.chTimeEventCmd:
			tet.handleCmd(cmd)
		case <-tet.timer.C:
			tet.triggerEvent()
			tet.refresh()
		}
	}
}

func (tet *TimeEventTrigger) refresh() {
	te := tet.findNearestEvent()
	if te == nil {
		tet.timer.Stop()
		return
	}
	tet.nearest = te
	tet.timer.Reset(te.delay)
}

func (tet *TimeEventTrigger) triggerEvent() {
	if tet.nearest == nil {
		__deTrigger__("[WARN]: trigger event but nearest is nil")
		return
	}
	tet.chTimeEventOut <- tet.nearest
	//del nearest te
	for i, v := range tet.tes {
		if v.Equals(tet.nearest) {
			__deTrigger__("te trigger, remove")
			tet.tes = append(tet.tes[:i], tet.tes[i+1:]...)
		}
	}
}

func (tet *TimeEventTrigger) handleCmd(cmd *TimeEventCmd) {
	switch cmd.tt {
	case TimeEventCmdTypeAdd:
		tet.addTimeEvent(cmd.te)
	}
	tet.refresh()
}

var ErrKeyDup = errors.New("time event key duplicate")

func (tet *TimeEventTrigger) addTimeEvent(te *TimeEvent) error {
	if tet.findTimeEvent(te.key) != nil {
		return ErrKeyDup
	}
	tet.tes = append(tet.tes, te)
	tet.refresh()
	return nil
}

func (tet *TimeEventTrigger) findTimeEvent(key string) *TimeEvent {
	for _, v := range tet.tes {
		if v.key == key {
			return v
		}
	}
	return nil
}

func (tet *TimeEventTrigger) findNearestEvent() *TimeEvent {
	var near *TimeEvent
	for _, te := range tet.tes {
		if near == nil {
			near = te
			continue
		}
		if near.delay > te.delay {
			near = te
		}
	}
	return near
}

const pathPrefix = "time_event"
const KEY_NAME = "key"
const DATA_NAME = "data"
const DELAY_NAME = "delay"

var (
	pathAdd = fmt.Sprintf("/%s/%s", pathPrefix, "add")
)

//ADD MOD DEL from http
func (tet *TimeEventTrigger) Post() {
	var reqObj struct {
		key   string `json:"key"`
		data  []byte `json:"data"`
		delay int64  `json:"delay"`
	}
	err := json.Unmarshal(tet.Ctx.Input.RequestBody, &reqObj)
	__deTrigger__("got json req obj :%v", reqObj)
	if err != nil {
		tet.Ctx.WriteString(fmt.Sprintf("wrong body parse json:%s", err))
		return
	}

	if reqObj.key == "" {
		tet.Ctx.WriteString(REP_NO_KEY)
		return
	}
	//delay is seconds
	if reqObj.delay < 0 {
		tet.Ctx.WriteString(REP_DELAY_WRONG)
		return
	}
	err = tet.addTimeEvent(&TimeEvent{
		key:   reqObj.key,
		data:  reqObj.data,
		delay: time.Second * time.Duration(reqObj.delay),
	})
	if err != nil {
		tet.Ctx.WriteString(fmt.Sprintf("err when add to time events:%s", err))
		return
	}
	tet.Ctx.WriteString(REP_OK)
	return
}

var (
	REP_ROUTER_NOT_FOUND = "404 router not found\n"
	REP_OK               = "OK\n"
	REP_DELAY_WRONG      = "param delay must be posive(in seconds)\n"
	REP_NO_KEY           = "no key\n"
	REP_BAD_REQ          = "request body or query error\n"
)
