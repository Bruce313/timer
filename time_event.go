package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
)

type TimeEventGenerator interface {
	getNext() *time.Duration
	genTimeEvent() *TimeEvent
}

type timerTimeEventGenerator struct {
	meta         *TimeEventMeta
	time2Trigger time.Time
	isTriggerd   bool
}

func newTimerTimeEventGeneratro(delay time.Duration, m *TimeEventMeta) *timerTimeEventGenerator {
	now := time.Now()
	return &timerTimeEventGenerator{
		meta:         m,
		time2Trigger: now.Add(delay),
	}
}

func (tteg *timerTimeEventGenerator) getNext() time.Duration {
	if tteg.isTriggerd {
		return -1
	}
	now := time.Now()
	return tteg.time2Trigger.Sub(now)
}

func (tteg *timerTimeEventGenerator) genTimeEvent() *TimeEvent {
	tteg.isTriggerd = true
	now := time.Now()
	return &TimeEvent{
		TimeEventMeta: tteg.meta,
		timeTriggered: &now,
	}
}

type intervalTimeEventGenerator struct {
	schedule   cron.Schedule
	lastTigger *time.Time
	meta       *TimeEventMeta
}

func newIntervalTimeEventGenerator(spec string, m *TimeEventMeta) (*intervalTimeEventGenerator, error) {
	schedule, err := cron.Parse(spec)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &intervalTimeEventGenerator{
		schedule:   schedule,
		lastTigger: &now,
		meta:       m,
	}, nil
}

func (iteg *intervalTimeEventGenerator) getNext() time.Duration {
	return time.Duration(iteg.schedule.Next(*iteg.lastTigger).UnixNano())
}

func (iteg *intervalTimeEventGenerator) genTimeEvent() *TimeEvent {
	now := time.Now()
	iteg.lastTigger = &now
	return &TimeEvent{
		TimeEventMeta: iteg.meta,
		timeTriggered: &now,
	}
}

type TimeEventMeta struct {
	key          string
	data         []byte
	timeRegister *time.Time
}

//TimeEvent contain key and data of time
type TimeEvent struct {
	*TimeEventMeta
	timeTriggered *time.Time
}

func (te *TimeEvent) String() string {
	return fmt.Sprintf("\n*****************\n key: %s,\n data: %s,\n time register: %s,\n**************\n",
		te.key, string(te.data), te.timeRegister.Format("2006-01-02 15:04:05"))
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
