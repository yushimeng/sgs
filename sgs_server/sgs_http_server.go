package sgs_server

import (
	"fmt"
	"net/http"
	"sgs/sgs_conf"
)

type httpHandleMap map[string]func(w http.ResponseWriter, r *http.Request)

type SgsHttpServer struct {
	Config    sgs_conf.ConfigHttpServer
	Server    http.Server
	HandleMgr httpHandleMap
}

func NewSgsHttpServer(cfg sgs_conf.ConfigHttpServer) *SgsHttpServer {
	s := new(SgsHttpServer)
	s.Config = cfg
	return s
}

func (s *SgsHttpServer) init() {
	s.HandleMgr = make(httpHandleMap)
	fmt.Println("SgsHttpServer init run...")
	s.HandleMgr["/status"] = HttpStatusHandle
	s.HandleMgr["/hi"] = HttpSayHiHandle
}

func (s *SgsHttpServer) Start() (err error) {
	s.init()
	addr := fmt.Sprintf(":%d", s.Config.HttpPort)

	for suffix, handle := range s.HandleMgr {
		http.HandleFunc(suffix, handle)
	}

	err = http.ListenAndServe(addr, nil)

	return err
}

func HttpStatusHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is status")
}

func HttpSayHiHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is hi")
}
