package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tj/go-debug"
)

var (
	deTrigger = debug.Debug("timer:trigger")
)

//TimeEventTrigger is a wrapper which triggers timeevent
type TimeEventTrigger struct {
	chTimeEventOut chan<- *TimeEvent
	chTimeEventCmd <-chan *TimeEventCmd
	//TODO heap it
	tes     []*TimeEvent
	timer   *time.Timer
	nearest *TimeEvent
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
		deTrigger("[WARN]: trigger event but nearest is nil")
		return
	}
	tet.chTimeEventOut <- tet.nearest
	//del nearest te
	for i, v := range tet.tes {
		if v.Equals(tet.nearest) {
			deTrigger("te trigger, remove")
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
func (tet *TimeEventTrigger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	deTrigger("got http req, path:%s", path)
	errParse := r.ParseForm()
	if errParse != nil {
		w.Write(REP_BAD_REQ)
		return
	}
	form := r.Form
	if path == pathAdd {
		deTrigger("match add, data:%s", form)
		key := form.Get(KEY_NAME)
		if key == "" {
			w.Write(REP_NO_KEY)
			return
		}
		data := form.Get(DATA_NAME)
		delay := form.Get(DELAY_NAME)
		//delay is seconds
		seconds, err := strconv.Atoi(delay)
		if err != nil || seconds < 0 {
			w.Write(REP_DELAY_WRONG)
			return
		}
		err = tet.addTimeEvent(&TimeEvent{
			key:   key,
			data:  []byte(data),
			delay: time.Second * time.Duration(seconds),
		})
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(REP_OK)
		return
	}
	w.Write(REP_ROUTER_NOT_FOUND)
}

var (
	REP_ROUTER_NOT_FOUND = []byte("404 router not found\n")
	REP_OK               = []byte("OK\n")
	REP_DELAY_WRONG      = []byte("param delay must be posive(in seconds)\n")
	REP_NO_KEY           = []byte("no key\n")
	REP_BAD_REQ          = []byte("request body or query error\n")
)
