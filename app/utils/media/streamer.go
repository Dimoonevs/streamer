package media

import (
	"flag"
	"fmt"
	_ "github.com/Dimoonevs/hls-parser/pkg/creator/hls"
	"github.com/Dimoonevs/hls-parser/pkg/domain"
	"log"
	"sync"
	"time"
)

var (
	countSegmentStream = 5
	publicURLMPL       = flag.String("publicURLPlaylist", "", "link public playlist")
)

type MediaStreamer struct {
	wg             *sync.WaitGroup
	seqNo          int
	Segments       map[string]*domain.MediaPlaylist
	MediaPlaylists []domain.MediaPlaylistForCreate
	mx             *sync.Mutex
}

func InitMediaStreamer() *MediaStreamer {
	return &MediaStreamer{
		wg:             &sync.WaitGroup{},
		seqNo:          0,
		Segments:       make(map[string]*domain.MediaPlaylist),
		mx:             &sync.Mutex{},
		MediaPlaylists: []domain.MediaPlaylistForCreate{},
	}
}

func (ms *MediaStreamer) CreateStreamMPL(playlists []*domain.MediaPlaylist) {
	go ms.createLink()
	if ms == nil {
		log.Fatal("MediaStreamer is nil")
	}
	for _, playlist := range playlists {
		ms.wg.Add(1)
		go ms.slideAndStream(playlist)
	}
	ms.wg.Wait()
}

func (ms *MediaStreamer) slideAndStream(playlist *domain.MediaPlaylist) {
	defer ms.wg.Done()
	if playlist == nil {
		log.Fatal("playlist is nil")
	}
	ms.mx.Lock()
	playlistStream, found := ms.Segments[playlist.Size]
	ms.mx.Unlock()
	if !found {
		playlistStream = &domain.MediaPlaylist{
			Size:           playlist.Size,
			TargetDuration: playlist.TargetDuration,
			SeqNo:          0,
			Segments:       []domain.Segment{},
		}
	}

	totalSegments := len(playlist.Segments)
	if totalSegments == 0 {
		return
	}

	if len(playlist.Segments) < countSegmentStream && len(playlistStream.Segments) < countSegmentStream {
		for _, segment := range playlist.Segments {
			playlistStream.Segments = append(playlistStream.Segments, segment)
			ms.mx.Lock()
			ms.Segments[playlist.Size] = playlistStream
			ms.mx.Unlock()
		}
		return
	}

	for i := 0; i < totalSegments; i++ {
		if len(playlistStream.Segments) >= countSegmentStream {
			playlistStream.Segments = playlistStream.Segments[1:]
		}
		playlistStream.Segments = append(playlistStream.Segments, playlist.Segments[i])

		if len(playlistStream.Segments) < countSegmentStream {
			continue
		}
		ms.mx.Lock()
		ms.Segments[playlist.Size] = playlistStream
		ms.mx.Unlock()

		totalDuration := 0
		for _, segment := range playlistStream.Segments {
			totalDuration += int(segment.Duration)
		}

		//streamPL := hls.GenerateMediaPlaylist(playlistStream)
		//log.Println(streamPL)

		time.Sleep(time.Duration(totalDuration/countSegmentStream) * time.Second)

		playlistStream.SeqNo++
	}
}

func (ms *MediaStreamer) createLink() {
	ms.MediaPlaylists = []domain.MediaPlaylistForCreate{}
	for _, mpl := range ms.Segments {
		ms.MediaPlaylists = append(ms.MediaPlaylists, domain.MediaPlaylistForCreate{
			MediaURL: fmt.Sprintf("%s/%s/playlist.m3u8", *publicURLMPL, mpl.Size),
			Size:     mpl.Size,
		})
	}

}
