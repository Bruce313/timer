package main

import (
	"errors"
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
	tes     []TimeEventGenerator
	timer   *time.Timer
	nearest TimeEventGenerator
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
		tes:            make([]TimeEventGenerator, 0),
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
			__deTrigger__("read time event cmd:%v", cmd)
			tet.handleCmd(cmd)
		case <-tet.timer.C:
			tet.triggerEvent()
			tet.refresh()
		}
	}
}

func (tet *TimeEventTrigger) refresh() {
	teg := tet.findNearestEventGenerator()
	if teg == nil {
		tet.timer.Stop()
		return
	}
	tet.nearest = teg
	tet.timer.Reset(teg.getNext())
}

func (tet *TimeEventTrigger) triggerEvent() {
	if tet.nearest == nil {
		__deTrigger__("[WARN]: trigger event but nearest is nil")
		return
	}
	te := tet.nearest.genTimeEvent()
	now := time.Now()
	te.timeTriggered = &now
	//TODO: deal with blocking
	tet.chTimeEventOut <- te
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
		tet.AddTimeEventGenerator(cmd.teg)
	}
	tet.refresh()
}

var ErrKeyDup = errors.New("time event key duplicate")

func (tet *TimeEventTrigger) AddTimeEventGenerator(teg TimeEventGenerator) error {
	for _, v := range tet.tes {
		if v.Equals(teg) {
			return ErrKeyDup
		}
	}
	tet.tes = append(tet.tes, teg)
	tet.refresh()
	return nil
}

func (tet *TimeEventTrigger) findNearestEventGenerator() TimeEventGenerator {
	var near TimeEventGenerator
	for _, te := range tet.tes {
		if near == nil {
			near = te
			continue
		}
		if near.getNext() > te.getNext() {
			near = te
		}
	}
	return near
}
