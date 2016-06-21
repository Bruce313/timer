package main

import (
	"fmt"
	"time"
)

//TimeEvent contain key and data of time
type TimeEvent struct {
	key           string        `json:"key"`
	data          []byte        `json:"data"`
	delay         time.Duration `json:"delay"`
	timeCreate    *time.Time    `json:"timeCreate"`
	timeTriggered *time.Time    `json:"timeTriggered"`
}

func (te *TimeEvent) String() string {
	return fmt.Sprintf("\n*****************\n key: %s,\n data: %s,\n delay: %d,\n**************\n",
		te.key, string(te.data), te.delay)
}

func (te *TimeEvent) Equals(ot *TimeEvent) bool {
	return te.key == ot.key
}

//TimeEventCmdType is types of operation to TimeEventCmd
type TimeEventCmdType int

const (
	_ TimeEventCmdType = iota
	//TimeEventCmdTypeAdd add
	TimeEventCmdTypeAdd
	//TimeEventCmdTypeMod mod
	TimeEventCmdTypeMod
	//TimeEventCmdTypeDel del
	TimeEventCmdTypeDel
)

//TimeEventCmd contain operation and config of timeevent
type TimeEventCmd struct {
	tt TimeEventCmdType
	te *TimeEvent
}
