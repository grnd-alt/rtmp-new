package hlsdata

import (
	"fmt"
	"strings"
)

type Playlist interface {
	GetString() string
}

type MediaPlaylist struct {
	lines          []string
	tsFileNames    []string
	targetDuration int
	sequenceNum    int
	maxFiles       int
}

func NewMediaPlaylist(targetDuration int, maxFiles int) *MediaPlaylist {
	return &MediaPlaylist{
		lines:          []string{},
		tsFileNames:    []string{},
		targetDuration: targetDuration,
		sequenceNum:    0,
		maxFiles:       maxFiles,
	}
}

func (playlist *MediaPlaylist) AddTsFile(newFileName string, duration float64) {
	if len(playlist.tsFileNames) == playlist.maxFiles {
		playlist.tsFileNames = playlist.tsFileNames[1:]
	}
	playlist.tsFileNames = append(playlist.tsFileNames, newFileName)
	playlist.lines = []string{
		"#EXTM3U",
		"#EXT-X-VERSION:3",
		fmt.Sprintf("#EXT-X-TARGETDURATION:%d", playlist.targetDuration),
		func() string {
			if playlist.sequenceNum > playlist.maxFiles {
				return fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%d", playlist.sequenceNum-playlist.maxFiles+2)
			} else {
				return fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%d", 5)
			}
		}(),
	}
	for _, fileName := range playlist.tsFileNames {
		playlist.lines = append(playlist.lines, fmt.Sprintf("#EXTINF:%f,", duration))
		playlist.lines = append(playlist.lines, fileName)
	}
	playlist.sequenceNum++
}

func (playlist *MediaPlaylist) GetString() string {
	return strings.Join(playlist.lines, "\n")
}

type MasterPlaylist struct {
	lines     []string
	playlists []MediaPlaylist
}

func NewMasterPlaylist(streamId string) *MasterPlaylist {
	return &MasterPlaylist{
		lines: []string{
			"#EXTM3U",
			// "#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID=\"chunked\",NAME=\"1080p60 (source)\",AUTOSELECT=YES,DEFAULT=YES",
			"#EXT-X-STREAM-INF:BANDWIDTH=128000,RESOLUTION=1920x1080,FRAME-RATE=60.000,CODEC=\"avc1.42e01e\"",
			fmt.Sprintf("/%s/1080p60/.m3u8", streamId),
		},
	}
}

func (playlist *MasterPlaylist) GetString() string {
	return strings.Join(playlist.lines, "\n")
}
