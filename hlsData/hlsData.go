package hlsdata

import (
	"fmt"
	"os"
	"time"

	"github.com/yapingcat/gomedia/go-codec"
)

type TransCoder struct {
	streams map[string]*StreamData
}

func NewTransCoder() TransCoder {
	transCoder := TransCoder{streams: map[string]*StreamData{}}
	return transCoder
}

func (t *TransCoder) GetStream(streamId string) (*StreamData, error) {
	stream, found := t.streams[streamId]
	if !found {
		return nil, fmt.Errorf("StreamId %s not found", streamId)
	}
	return stream, nil
}

func (t *TransCoder) OnPublish(streamId string) {
	fmt.Println("PUBLISHED STREAM")
	t.streams[streamId] = NewStreamData(streamId, 6)
	os.RemoveAll(fmt.Sprintf("./%s", streamId))
}

func (t *TransCoder) OnFrame(codecid int, frame []byte, streamId string) {
	stream, _ := t.GetStream(streamId)
	if codecid == int(codec.CODECID_VIDEO_H264) {
		stream.videoData = append(stream.videoData, frame...)
	} else if codecid == int(codec.CODECID_AUDIO_AAC) {
		stream.audioData = append(stream.audioData, frame...)
	} else {
		panic("UNKNOWN CODEC ID")
	}
	if time.Since(stream.lastMux) > 2*time.Second {
		go stream.addTSFile(stream.audioData, stream.videoData, time.Since(stream.lastMux).Seconds())
		stream.audioData = []byte{}
		stream.videoData = []byte{}
		stream.lastMux = time.Now()
	}
}

func (t *TransCoder) GetPlaylist(streamId string, playlistName string) string {
	stream, err := t.GetStream(streamId)
	if err != nil {
		return ""
	}
	return stream.Playlists[0].GetString()
}

func (t *TransCoder) GetMasterPlaylist(streamId string) string {
	stream, err := t.GetStream(streamId)
	if err != nil {
		return ""
	}
	// return stream.Playlists[0].GetString()
	return stream.MasterPlaylist.GetString()
}
