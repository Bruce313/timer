package main

import (
	. "github.com/tj/go-debug"
)

var __deRegCli__ = Debug("timer:reg_client")

type RegisterClient struct {
	name string
}

func NewRegisterClient(name string) *RegisterClient {
	return &RegisterClient{
		name: name,
	}
}

func (rc *RegisterClient) deliver(te *TimeEvent, name string) error {
	__deRegCli__("client for name:%s, try to deliver timeevent:%s", name, te)
	return nil
}
