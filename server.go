package main

import (
	"net/http"

	. "github.com/tj/go-debug"
)

var __deMain__ = Debug("timer:main")

func main() {
	chTimeEventOut := make(chan *TimeEvent)
	chTimeEventCmd := make(chan *TimeEventCmd, 10)
	trigger := NewTimeEventTrigger(chTimeEventOut, chTimeEventCmd)
	go trigger.Begin()
	debugTecm := NewDebugTECM()
	regMgrTecm := NewRegisterMgr()
	go broadcastEvent(chTimeEventOut, debugTecm, regMgrTecm)
	//http register client
	http.Handle("/time_event", trigger)
	http.Handle("/handler/http", regMgrTecm)
	http.ListenAndServe(":6200", nil)
}

func broadcastEvent(chTE <-chan *TimeEvent, heads ...TimeEventHandler) {
	for {
		__deMain__("wait for te")
		te := <-chTE
		for _, head := range heads {
			err := head.HandleTimeEvent(te)
			if err != nil {
				__deMain__("handle error:%s, quit", err)
				break
			}
		}
	}
}
