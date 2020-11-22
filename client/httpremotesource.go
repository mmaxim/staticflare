package client

import (
	"io/ioutil"
	"net/http"

	"mmaxim.org/staticflare/common"
)

type HTTPResponseHandler interface {
	ParseRemoteIP(resp string) (string, error)
}

type HTTPRemoteSource struct {
	*common.DebugLabeler
	url         string
	respHandler HTTPResponseHandler
}

func NewHTTPRemoteSource(url string, respHandler HTTPResponseHandler) *HTTPRemoteSource {
	return &HTTPRemoteSource{
		DebugLabeler: common.NewDebugLabeler("HTTPRemoteSource"),
		url:          url,
		respHandler:  respHandler,
	}
}

func (s *HTTPRemoteSource) GetRemoteIP() (res string, err error) {
	r, err := http.Get(s.url)
	if err != nil {
		s.Debug("GetRemoteIP: failed to make HTTP req: %s", err)
		return res, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.Debug("GetRemoteIP: failed to read response: %s", err)
		return res, err
	}
	if res, err = s.respHandler.ParseRemoteIP(string(body[:])); err != nil {
		s.Debug("GetRemoteIP: failed to handle response: %s", err)
	}
	return res, err
}
