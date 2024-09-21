package hlsServer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type (
	GetMasterPlaylist func(streamId string) string
	GetSegment        func(streamID string, segment string) string
	GetPlaylist       func(streamId string, playlist string) string
	HlsServer         struct {
		Hostname          string `json:"hostname"`
		Port              int    `json:"port"`
		mux               http.Handler
		getMasterPlaylist GetMasterPlaylist
		getPlaylist       GetPlaylist
		getSegment        GetSegment
	}
)

func (server *HlsServer) OnGetMasterPlaylist(handler GetMasterPlaylist) {
	server.getMasterPlaylist = handler
}

func (server *HlsServer) OnGetPlaylist(handler GetPlaylist) {
	server.getPlaylist = handler
}

func (server *HlsServer) OnGetSegment(handler GetSegment) {
	server.getSegment = handler
}

func handleHomeRequest(w http.ResponseWriter, req *http.Request) {
	data := HlsServer{
		Hostname: "test",
		Port:     192134,
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (server *HlsServer) handleMasterPlaylistRequest(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, server.getMasterPlaylist(req.PathValue("streamId")))
}

func (server *HlsServer) handlePlaylistRequest(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, server.getPlaylist(req.PathValue("streamId"), req.PathValue("playlist")))
}

func (server *HlsServer) handleSegmentRequest(w http.ResponseWriter, req *http.Request) {
	if strings.HasSuffix(req.PathValue("segment"), ".mp4") {
		w.Header().Set("Content-Type", "video/mp4")
	}

	fmt.Fprint(w, server.getSegment(req.PathValue("streamId"), req.PathValue("segment")))
}

func corsMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Warn("Request")
		w.Header().Set("Access-control-allow-origin", "*")
		next.ServeHTTP(w, r)
	})
}

func InitHlsServer(hostname string, port int) *HlsServer {
	server := HlsServer{Hostname: hostname, Port: port}
	mux := http.NewServeMux()

	mux.HandleFunc("/{streamId}/.m3u8", server.handleMasterPlaylistRequest)
	// mux.HandleFunc("/{streamId}/", server.handleSegmentRequest)
	mux.HandleFunc("/{streamId}/{playlist}/.m3u8", server.handlePlaylistRequest)
	mux.HandleFunc("/{streamId}/{playlist}/{segment}", server.handleSegmentRequest)
	mux.HandleFunc("/", handleHomeRequest)
	wrappedMux := corsMiddleWare(mux)
	server.mux = wrappedMux
	return &server
}

func (server *HlsServer) StartHlsServer() {
	http.ListenAndServe(fmt.Sprintf(":%d", server.Port), server.mux)
}
