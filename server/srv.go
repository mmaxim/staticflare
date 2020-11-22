package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"mmaxim.org/staticflare/common"
)

type Server struct {
	*common.DebugLabeler
	addr string
}

func NewServer(addr string) *Server {
	return &Server{
		DebugLabeler: common.NewDebugLabeler("Server"),
		addr:         addr,
	}
}

func (s *Server) handleInfo(w http.ResponseWriter, req *http.Request) {
	ip := strings.Split(req.RemoteAddr, ":")[0]
	ret := common.InfoResponse{
		RemoteIP: ip,
	}
	b, err := json.Marshal(ret)
	if err != nil {
		s.Debug("handleInfo: failed to marshal response: %s", err)
		return
	}
	if _, err := w.Write(b); err != nil {
		s.Debug("handleInfo: failed to write response: %s", err)
		return
	}
	s.Debug("handleInfo(): ret: %s", string(b))
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/info", s.handleInfo)
	srv := http.Server{
		Addr:    s.addr,
		Handler: mux,
	}
	s.Debug("starting up on: %s", s.addr)
	return srv.ListenAndServe()
}
