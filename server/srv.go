package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) handleInfo(w http.ResponseWriter, req *http.Request) {
	ip := strings.Split(req.RemoteAddr, ":")[0]
	ret := infoResponse{
		RemoteIP: ip,
	}
	b, err := json.Marshal(ret)
	if err != nil {
		log.Printf("handleInfo: failed to marshal response: %s\n", err)
		return
	}
	if _, err := w.Write(b); err != nil {
		log.Printf("handleInfo: failed to write response: %s\n", err)
		return
	}
	log.Printf("handleInfo(): ret: %s\n", string(b))
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/info", s.handleInfo)
	srv := http.Server{
		Addr:    s.addr,
		Handler: mux,
	}
	log.Printf("starting up on: %s\n", s.addr)
	return srv.ListenAndServe()
}
