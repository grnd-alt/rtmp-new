package main

import (
	"fmt"
	"os"
	hlsdata "rtmp-new/m/v2/hlsData"
	"rtmp-new/m/v2/hlsServer"

	rtmp "github.com/grnd-alt/rtmpServer"
	"github.com/yapingcat/gomedia/go-codec"
)

type fileWriter struct {
	videoFile *os.File
	audioFile *os.File
}

func (writer *fileWriter) onPublish(streamId string) {
	var err error
	writer.videoFile, err = os.OpenFile(fmt.Sprintf("%s.h264", streamId), os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("cannot create File to write to")
	}
	writer.audioFile, err = os.OpenFile(fmt.Sprintf("%s.aac", streamId), os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("Cannot create File to write to")
	}
	fmt.Printf("PUBLISH: %s\n", streamId)
}

func (writer *fileWriter) onFrame(codecid int, frame []byte, streamId string) {
	if codecid == int(codec.CODECID_VIDEO_H264) {
		writer.videoFile.Write(frame)
	} else if codecid == int(codec.CODECID_AUDIO_AAC) {
		writer.audioFile.Write(frame)
	} else {
		panic("UNKNOWN CODEC ID")
	}
	// fmt.Printf("FRAME: %s\n", streamId)
}

func main() {
	// writer := fileWriter{}
	transCoder := hlsdata.NewTransCoder()
	go rtmp.StartServer("0.0.0.0", 1935, transCoder.OnFrame, transCoder.OnPublish)
	HlsServer := hlsServer.InitHlsServer("0.0.0.0", 8080)
	HlsServer.OnGetMasterPlaylist(transCoder.GetMasterPlaylist)
	HlsServer.OnGetPlaylist(transCoder.GetPlaylist)

	HlsServer.OnGetSegment(func(streamID, segment string) string {
		content, err := os.ReadFile(fmt.Sprintf("./%s/%s", streamID, segment))
		if err != nil {
			panic(err)
		}
		return string(content)
	})
	HlsServer.StartHlsServer()
	fmt.Println("done lol")
}
