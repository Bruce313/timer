package main

import (
    "time"
    "fmt"
)

//TimeEvent contain key and data of time
type TimeEvent struct {
    key string
    data []byte
    delay time.Duration
}

func (te *TimeEvent) String() string {
    return fmt.Sprintf("\n*****************\n key: %s,\n data: %s,\n delay: %d,\n**************\n",
        te.key, string(te.data), te.delay)
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
