package main

import (
	"fmt"
	"os"
	"rtmp-new/m/v2/hlsServer"
	"rtmp-new/m/v2/rtmp"

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
	fmt.Printf("FRAME: %s\n", streamId)
}

type transCoder struct {
	videoData []byte
	audioData []byte
}

func (transcoder *transCoder) onFrame(codecid int, frame []byte, streamId string) {
	if codecid == int(codec.CODECID_VIDEO_H264) {
		transcoder.videoData = append(transcoder.videoData, frame...)
	} else if codecid == int(codec.CODECID_AUDIO_AAC) {
		transcoder.audioData = append(transcoder.audioData, frame...)
	} else {
		panic("UNKNOWN CODEC ID")
	}
	fmt.Printf("FRAME: %s\n", streamId)
}
func main() {
	writer := fileWriter{}
	go rtmp.StartServer("localhost", 1935, writer.onFrame, writer.onPublish)
	HlsServer := hlsServer.InitHlsServer("localhost", 8080)
	HlsServer.OnGetMasterPlaylist(func(streamId string) string {
		content, err := os.ReadFile(fmt.Sprintf("./%s/master.m3u8", streamId))
		if err != nil {
			panic(err)
		}
		return string(content)
	})

	HlsServer.OnGetSegment(func(streamID, segment string) string {
		content, err := os.ReadFile(fmt.Sprintf("./%s/%s", streamID, segment))
		if err != nil {
			panic(err)
		}
		return string(content)
	})
	HlsServer.StartHlsServer()
}
