package main

type RegisterClientTCP struct {
	RegisterClientMatcher
	name string
}

func NewRegisterClientTCP(m RegisterClientMatcher, name string) RegisterClient {
	return &RegisterClientTCP{
		name: name,
		RegisterClientMatcher: m,
	}
}

func (rc *RegisterClientTCP) Deliver(te *TimeEvent) error {
	__deRegCli__("client for name:%s, try to deliver timeevent:%s", rc.name, te)
	return nil
}
