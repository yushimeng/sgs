package sgs_server

import (
	"fmt"
	"net/http"
	"sgs/sgs_conf"
	"sgs/util"
)

type httpHandleMap map[string]func(w http.ResponseWriter, r *http.Request)

type SgsHttpServer struct {
	server    *SgsServer
	Config    sgs_conf.ConfigHttpServer
	Server    http.Server
	HandleMgr httpHandleMap
}

func NewSgsHttpServer(svr *SgsServer, cfg sgs_conf.ConfigHttpServer) *SgsHttpServer {
	httpServer := new(SgsHttpServer)
	httpServer.server = svr
	httpServer.Config = cfg
	return httpServer
}

func (httpServer *SgsHttpServer) init() {
	httpServer.HandleMgr = make(httpHandleMap)
	fmt.Println("SgsHttpServer init set http url handle...")
	httpServer.HandleMgr["/status"] = HttpStatusHandle
	httpServer.HandleMgr["/hi"] = HttpSayHiHandle
}

func (httpServer *SgsHttpServer) Start() {
	httpServer.init()
	addr := fmt.Sprintf(":%d", httpServer.Config.HttpPort)

	for suffix, handle := range httpServer.HandleMgr {
		http.HandleFunc(suffix, handle)
	}

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		util.AbnormalExit()
	}
}

func HttpStatusHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is status")
}

func HttpSayHiHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is hi")
}
