package hlsdata

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/yapingcat/gomedia/go-codec"
	"github.com/yapingcat/gomedia/go-mpeg2"
)

type StreamData struct {
	streamId       string
	muxer          *mpeg2.TSMuxer
	targetDuration uint
	videoData      []byte
	audioData      []byte
	lastMux        time.Time
	segmentCount   uint
	videoPid       uint16
	audioPid       uint16
	Playlists      []Playlist
	MasterPlaylist MasterPlaylist
	pts            uint64
	dts            uint64
}

func NewStreamData(streamId string, targetDuration uint) *StreamData {
	stream := &StreamData{
		streamId:       streamId,
		muxer:          mpeg2.NewTSMuxer(),
		targetDuration: targetDuration,
		videoData:      []byte{},
		audioData:      []byte{},
		lastMux:        time.Now(),
		segmentCount:   0,
		videoPid:       0,
		audioPid:       0,
		Playlists: []Playlist{
			NewMediaPlaylist(int(targetDuration), 10),
		},
		MasterPlaylist: *NewMasterPlaylist(streamId),
	}
	err := os.MkdirAll(fmt.Sprintf("./%s", stream.streamId), os.ModePerm)
	if err != nil {
		panic(err)
	}
	return stream
}

func (s *StreamData) writeVideo(wg *sync.WaitGroup, data []byte) {
	defer wg.Done()
	codec.SplitFrameWithStartCode(data, func(nalu []byte) bool {
		if codec.H264NaluType(nalu) <= codec.H264_NAL_I_SLICE {
			s.pts += 33
			s.dts += 33
		}
		s.muxer.Write(s.videoPid, nalu, s.pts, s.dts)
		return true
	})
}

func (s *StreamData) writeAudio(wg *sync.WaitGroup, data []byte) {
	defer wg.Done()
	// i := 0
	codec.SplitAACFrame(data, func(aac []byte) {
		// if i < 3 {
		// s.pts += 23
		// s.dts += 23
		// i++
		// } else {
		s.pts += 21
		s.dts += 21
		// i = 0
		// }
		s.muxer.Write(s.audioPid, aac, s.pts, s.dts)
	})
}

func (s *StreamData) addTSFile(audio []byte, video []byte, duration float64) {
	s.videoPid = 0
	s.muxer = mpeg2.NewTSMuxer()
	if s.videoPid == 0 {
		s.videoPid = s.muxer.AddStream(mpeg2.TS_STREAM_H264)
		s.audioPid = s.muxer.AddStream(mpeg2.TS_STREAM_AAC)
	}

	err := os.MkdirAll(fmt.Sprintf("./%s", s.streamId), os.ModePerm)
	if err != nil {
		panic(err)
	}
	if s.segmentCount > 20 {
		os.Remove(fmt.Sprintf("./%s/%d.ts", s.streamId, s.segmentCount-21))
	}
	tsFile, err := os.OpenFile(fmt.Sprintf("./%s/%d.ts", s.streamId, s.segmentCount), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer tsFile.Close()

	s.muxer.OnPacket = func(pkg []byte) {
		tsFile.Write(pkg)
	}
	var wg sync.WaitGroup
	s.pts = 0
	s.dts = 0
	wg.Add(1)
	// go s.writeAudio(&wg, audio)
	go s.writeVideo(&wg, video)
	wg.Wait()
	s.Playlists[0].(*MediaPlaylist).AddTsFile(fmt.Sprintf("%d.ts", s.segmentCount), duration)
	s.segmentCount++
}
