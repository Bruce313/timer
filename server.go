package main

import (
	"github.com/astaxie/beego"
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
	go broadcastEvent(chTimeEventOut, debugTecm, regMgrTecm)
	//http register client
	beego.Router("/client", NewRegisterMgr())
	beego.Router("/time_event", trigger)
	beego.Run()
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
