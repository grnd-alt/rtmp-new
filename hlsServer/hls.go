package hlsServer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type GetMasterPlaylist func(streamId string) string
type GetSegment func(streamID string, segment string) string
type HlsServer struct {
	Hostname          string `json:"hostname"`
	Port              int    `json:"port"`
	mux               http.Handler
	getMasterPlaylist GetMasterPlaylist
	getSegment        GetSegment
}

func (server *HlsServer) OnGetMasterPlaylist(handler GetMasterPlaylist) {
	server.getMasterPlaylist = handler
}

func (server *HlsServer) OnGetSegment(handler GetSegment) {
	server.getSegment = handler
}

func getPlaylistFile() string {
	playlist1 := "#EXT-X-STREAM-INF:CODECS=\"avc1.64002a,mp4a.40.2\",RESOLUTION=1920x1080,FRAME-RATE=60.000\nprog_index.m3u8\n"
	return fmt.Sprintf("#EXTM3U\n#EXT-X-VERSION:6\n#EXT-X-INDEPENDENT-SEGMENTS\n%s", playlist1)
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

func (server *HlsServer) handleSegmentRequest(w http.ResponseWriter, req *http.Request) {

	if strings.HasSuffix(req.PathValue("segment"), ".mp4") {
		w.Header().Set("Content-Type", "video/mp4")
	}
	fmt.Fprint(w, server.getSegment(req.PathValue("streamId"), req.PathValue("segment")))
}

func handlePlaylistRequest(w http.ResponseWriter, req *http.Request) {
	header := "#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-TARGETDURATION:7\n#EXT-X-MEDIA-SEQUENCE:1\n"
	fmt.Fprintf(w, "%s#EXT-X-MAP:URI=\".mp4\",BYTERANGE=\"718@0\"\n#EXTINF:7.0,\nsegment/.mp4", header)
}

func handleMP4Request(w http.ResponseWriter, req *http.Request) {
	b, err := os.ReadFile("2024-08-09_11-21-20.mp4")
	if err != nil {
		fmt.Fprintf(w, "%s, %s", req.PathValue("streamId"), req.PathValue("segment"))
		return
	}
	w.Header().Set("Content-Type", "video/mp4")
	w.Write(b)
}

func corsMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("request: ", r.RequestURI)
		w.Header().Set("Access-control-allow-origin", "*")
		next.ServeHTTP(w, r)
	})
}

func InitHlsServer(hostname string, port int) *HlsServer {
	server := HlsServer{Hostname: hostname, Port: port}
	mux := http.NewServeMux()

	mux.HandleFunc("/stream/{streamId}/.m3u8", server.handleMasterPlaylistRequest)
	mux.HandleFunc("/stream/{streamId}/{segment}", server.handleSegmentRequest)
	mux.HandleFunc("/", handleHomeRequest)
	wrappedMux := corsMiddleWare(mux)
	server.mux = wrappedMux
	return &server
}

func (server *HlsServer) StartHlsServer() {
	http.ListenAndServe(fmt.Sprintf(":%d", server.Port), server.mux)
}
