package main

import (
	"fmt"
	"os"
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

func main() {
	writer := fileWriter{}
	rtmp.StartServer("localhost", 1935, writer.onFrame, writer.onPublish)
}
