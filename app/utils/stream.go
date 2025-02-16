package utils

import (
	hls2 "github.com/Dimoonevs/hls-parser/pkg/creator/hls"
	"github.com/Dimoonevs/hls-parser/pkg/domain"
	"github.com/Dimoonevs/hls-parser/pkg/parser/hls"
	"hls-streamer/app/models"
	"hls-streamer/app/utils/media"
	"log"
	"net/http"
	"strings"
	"sync"
)

type MasterPlaylist struct {
	URI string
}

type Stream struct {
	ms    *media.MediaStreamer
	mutex *sync.Mutex
}

func InitStream() *Stream {
	return &Stream{
		mutex: &sync.Mutex{},
		ms:    media.InitMediaStreamer(),
	}
}

func (s *Stream) StartStream(videoContent *models.Video) {
	wg := &sync.WaitGroup{}
	mpls := make([]*domain.MediaPlaylist, 0)

	for _, plstr := range videoContent.Formats {
		wg.Add(1)
		go func(plstr *models.VideoFormat) {
			defer wg.Done()
			mpl, err := hls.ParseMediaPlaylist(plstr.URL)
			mpl.Size = plstr.Resolution
			if err != nil {
				log.Fatal(err)
			}
			s.mutex.Lock()
			mpls = append(mpls, mpl)
			s.mutex.Unlock()
		}(plstr)
	}

	wg.Wait()

	if len(mpls) == 0 {
		log.Fatal("mpls is empty, nothing to stream")
	}

	s.ms.CreateStreamMPL(mpls)
}

func (s *Stream) HandlePlaylist(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid request path", http.StatusBadRequest)
		return
	}
	size := pathParts[1]

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	playlist, found := s.ms.Segments[size]
	if !found {
		http.Error(w, "Not found", http.StatusNotFound)
	}

	_, err := w.Write([]byte(hls2.GenerateMediaPlaylist(playlist)))
	if err != nil {
		log.Println("Ошибка при отправке плейлиста:", err)
		return
	}
}
