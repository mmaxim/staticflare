package client

import (
	"io/ioutil"
	"net/http"
	"time"

	"mmaxim.org/staticflare/common"
)

type HTTPResponseHandler interface {
	ParseRemoteIP(resp string) (string, error)
}

type HTTPRemoteSource struct {
	*common.DebugLabeler
	url         string
	respHandler HTTPResponseHandler
	stats       common.StatsProvider
}

func NewHTTPRemoteSource(url string, respHandler HTTPResponseHandler, stats common.StatsProvider) *HTTPRemoteSource {
	return &HTTPRemoteSource{
		DebugLabeler: common.NewDebugLabeler("HTTPRemoteSource"),
		url:          url,
		respHandler:  respHandler,
		stats:        stats.SetPrefix("HTTPRemoteSource"),
	}
}

func (s *HTTPRemoteSource) GetRemoteIP() (res string, err error) {
	start := time.Now()
	defer func() {
		if err == nil {
			s.stats.Value("GetRemoteIP - speed ms", float64(time.Since(start).Milliseconds()))
		}
	}()

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
