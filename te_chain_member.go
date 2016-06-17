package main

import (
	. "github.com/tj/go-debug"
)

type TimeEventChainMember struct {
	next *TimeEventChainMember
}

func NewTimeEventChainMember(n *TimeEventChainMember) *TimeEventChainMember {
	return &TimeEventChainMember{
		next: n,
	}
}

func (tecm *TimeEventChainMember) HandleTimeEvent(te *TimeEvent) error {
	if tecm.next != nil {
		return tecm.next.HandleTimeEvent(te)
	}
	return nil
}

func (tecm *TimeEventChainMember) AppendMember(tm *TimeEventChainMember) {
	tecm.next = tm
}

const DEBUG_NAME = "timer:TimeEventChainMemeber:Debug"
const CHAIN_BUF_SIZE = 10

type DebugTECM struct {
	*TimeEventChainMember
	debug DebugFunction
}

func NewDebugTECM() *DebugTECM {
	return &DebugTECM{
		TimeEventChainMember: NewTimeEventChainMember(nil),
		debug:                Debug(DEBUG_NAME),
	}
}

func (dtecm *DebugTECM) HandleTimeEvent(te *TimeEvent) error {
	dtecm.debug("收到TimeEvent:%s", te)
	return dtecm.TimeEventChainMember.HandleTimeEvent(te)
}

type TimeEventHandler interface {
	//error 代表需要停止的错误　业务上的错误写入te中
	HandleTimeEvent(*TimeEvent) error
}
