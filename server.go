package main

import (
	"fmt"
	"net/http"

	. "github.com/tj/go-debug"
)

var deMain = Debug("timer:main")

func main() {
	chTimeEventOut := make(chan *TimeEvent)
	chTimeEventCmd := make(chan *TimeEventCmd)
	trigger := NewTimeEventTrigger(chTimeEventOut, chTimeEventCmd)
	go trigger.Begin()
	debugTecm := NewDebugTECM()
	regMgrTecm := NewRegisterMgr()
	http.Handle(fmt.Sprintf("/%s/", pathPrefix), trigger)
	go broadcastEvent(chTimeEventOut, debugTecm, regMgrTecm)
	http.ListenAndServe(":6200", nil)
}

func broadcastEvent(chTE <-chan *TimeEvent, heads ...TimeEventHandler) {
	for {
		deMain("wait for te")
		te := <-chTE
		for _, head := range heads {
			err := head.HandleTimeEvent(te)
			if err != nil {
				deMain("handle error:%s, quit", err)
				break
			}
		}
	}
}
