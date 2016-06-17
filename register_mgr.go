package main

import (
	"regexp"
	"strings"

	. "github.com/tj/go-debug"
)

var __deRegMgr__ = Debug("timer:registerMgr")

type RegisterMgr struct {
	fixedKey2Regs     map[string][]string         //map from te.key -> register names
	pattern2Regs      map[*regexp.Regexp][]string //pattern regexp -> register names
	connectedClients  map[string]RegisterClient   //client.name -> client
	delayedTimeEvents []*TimeEventPublish
}

func NewRegisterMgr() *RegisterMgr {
	return &RegisterMgr{
		fixedKey2Regs:    make(map[string][]string),
		pattern2Regs:     make(map[*regexp.Regexp][]string),
		connectedClients: make(map[string]RegisterClient),
	}
}

func (rgm *RegisterMgr) HandleTimeEvent(te *TimeEvent) error {
	__deRegMgr__("got time event to handle:%s", te)
	key := te.key
	for k, v := range rgm.fixedKey2Regs {
		if k == key {
			rgm.Publish(te, v...)
		}
	}
	for r, v := range rgm.pattern2Regs {
		if r.MatchString(key) {
			rgm.Publish(te, v...)
		}
	}
	return nil
}

func (rgm *RegisterMgr) Publish(te *TimeEvent, names ...string) {
	__deRegMgr__("publish te:%s, to names:%s clients", te, strings.Join(names, ","))
	for _, n := range names {
		c, ok := rgm.connectedClients[n]
		if !ok {
			rgm.delay(te, n)
		}
		err := c.deliver(te, n)
		if err != nil {
			rgm.delay(te, n)
		}
		//TODO remove client if err is fatal(ie. client closed)
	}
}

func (rgm *RegisterMgr) delay(te *TimeEvent, name string) {
	__deRegMgr__("clients not found for te:%s with name %s. delay")
	rgm.delayedTimeEvents = append(rgm.delayedTimeEvents, &TimeEventPublish{
		name: name,
		te:   te,
	})
}

type TimeEventPublish struct {
	name string //register name
	te   *TimeEvent
}
