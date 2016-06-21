package main

import (
	"regexp"

	. "github.com/tj/go-debug"
)

var __deRegCli__ = Debug("timer:reg_client")

type RegisterClient interface {
	RegisterClientDeliver
	RegisterClientMatcher
}

type RegisterClientDeliver interface {
	Deliver(te *TimeEvent) error
}

type RegisterClientMatcher interface {
	Match(key string) bool
}

type FixedKeyMatcher struct {
	key string
}

func NewFixedKeyMatcher(key string) *FixedKeyMatcher {
	return &FixedKeyMatcher{
		key: key,
	}
}

func (fnm *FixedKeyMatcher) Match(key string) bool {
	return fnm.key == key
}

type RegexpKeyMatcher struct {
	reg *regexp.Regexp
}

func NewRegexpKeyMatcher(regString string) (*RegexpKeyMatcher, error) {
	r, err := regexp.Compile(regString)
	if err != nil {
		return nil, err
	}
	return &RegexpKeyMatcher{
		reg: r,
	}, nil
}

func (rnm *RegexpKeyMatcher) Match(key string) bool {
	return rnm.reg.MatchString(key)
}
