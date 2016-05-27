package main

import (
    "time"
    "github.com/tj/go-debug"
)

var (
    deTrigger = debug.Debug("trigger")
)

//TimeEventTrigger is a wrapper which triggers timeevent
type TimeEventTrigger struct {
   chTimeEventOut chan <- *TimeEvent
   chTimeEventCmd <- chan *TimeEventCmd
   //TODO heap it
   tes []*TimeEvent
   timer *time.Timer
   nearest *TimeEvent
}

//NewTimeEventTrigger create TimeEventTrigger
func NewTimeEventTrigger(ch chan <- *TimeEvent, chType <- chan *TimeEventCmd) *TimeEventTrigger {
    timer := time.NewTimer(1 * time.Hour)
    timer.Stop()
    return &TimeEventTrigger {
        chTimeEventOut: ch,
        chTimeEventCmd: chType,
        timer: timer,
        tes: make([]*TimeEvent, 0),
    }
}

// //Add append timeevent to a key
// func (tet *TimeEventTrigger)Add(te *TimeEvent) {
//     tet.tes = append(tet, te)
// }

//Begin start wait timeevent and trigger it
func (tet *TimeEventTrigger)Begin() {
    for {
        select {
        case cmd := <- tet.chTimeEventCmd:
                tet.handleCmd(cmd)
            case  <- tet.timer.C:
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

func (tet *TimeEventTrigger)triggerEvent() {
    if tet.nearest == nil {
        deTrigger("[WARN]: trigger event but nearest is nil")
        return
    }
    tet.chTimeEventOut <- tet.nearest
}

func (tet *TimeEventTrigger) handleCmd(cmd *TimeEventCmd) {
    switch cmd.tt {
    case TimeEventCmdTypeAdd:
        tet.addTimeEvent(cmd.te)
    }
    tet.refresh()
}

func (tet *TimeEventTrigger) addTimeEvent(te *TimeEvent) {
    tet.tes = append(tet.tes, te)
    tet.refresh()
}

func(tet *TimeEventTrigger) findNearestEvent() *TimeEvent {
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
