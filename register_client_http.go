package main

import (
	"bytes"
	"net/http"
	"strings"
)

//use http request to notify
type RegisterClientHTTP struct {
	RegisterClientMatcher
	url  string
	name string
}

//DP: bridge
func NewRegisterClientHTTP(m RegisterClientMatcher, name, url string) *RegisterClientHTTP {
	return &RegisterClientHTTP{
		RegisterClientMatcher: m,
		url:  url,
		name: name,
	}
}

const REPLACE_KEY = "${key}"

func (rch *RegisterClientHTTP) Deliver(te *TimeEvent) error {
	//replace url with key ${key}
	replaceUrl := strings.Replace(rch.url, REPLACE_KEY, te.key, -1)
	bodyBuf := rch.buildBody(te)
	bodyReader := bytes.NewReader(bodyBuf)
	__deRegCli__("deliver: body:%s, url:%s", []byte(bodyBuf), replaceUrl)
	_, err := http.Post(replaceUrl, "application/json", bodyReader)
	if err != nil {
		__deRegCli__("post err:%s", err)
		return err
	}
	return nil
}

func (rch *RegisterClientHTTP) buildBody(te *TimeEvent) []byte {
	return []byte(`OK`)
	//return []byte(fmt.Sprintf(`
	//	{
	//		"key": "%s",
	//		"data": "%s",
	//		"delay": %d,
	//		"name": "%s"
	//	}
	//`, te.key, te.data, te.delay, rch.name))
}
