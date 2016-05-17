package main

import (
    "time"
)

//TimeEventTrigger is a wrapper which triggers timeevent 
type TimeEventTrigger struct {
   chTimeEvent chan <- *TimeEvent 
   //TODO heap it
   tes []*TimeEvent
}

//NewTimeEventTrigger create TimeEventTrigger
func NewTimeEventTrigger(ch chan <- *TimeEvent) *TimeEventTrigger {
    return &TimeEventTrigger {
        chTimeEvent: ch,
    }
}

//Add append timeevent to a key
func (tet *TimeEventTrigger)Add(te *TimeEvent) {
    tet.tes = append(tet, te) 
}

//Begin start wait timeevent and trigger it
func 