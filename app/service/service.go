package service

import (
	"flag"
	"fmt"
	"hls-streamer/app/models"
	"hls-streamer/app/repo/mysql"
	"hls-streamer/app/utils"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	portStreamer   = flag.Int("portStreamer", 8080, "The port streamer endpoint")
	listVideo      []models.Video
	listVideoMutex sync.Mutex
)

type Segment struct {
	Duration float64
	URI      string
}

func Run() error {
	var err error
	listVideo, err = mysql.GetConnection().GetVideoContent()
	if err != nil {
		return err
	}
	log.Println("Gets List Video")

	stream := utils.InitStream()
	go func() {
		log.Printf("Start streamer on port :%d\n", *portStreamer)
		http.HandleFunc("/master.m3u8", stream.HandleMasterPlaylist)
		http.HandleFunc("/", stream.HandlePlaylist)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *portStreamer), nil))
	}()

	go updateVideoListPeriodically()

	for {
		listVideoMutex.Lock()
		for _, video := range listVideo {
			stream.StartStream(&video)
		}
		listVideoMutex.Unlock()
	}
}

func updateVideoListPeriodically() {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C

		log.Println("🔄 Start update Video")

		newList, err := mysql.GetConnection().GetVideoContent()
		if err != nil {
			log.Printf("❌ Error get videos: %v", err)
			continue
		}

		log.Println("Gets List Video")

		listVideo = newList

		log.Println("✅ List video Update successful")
	}
}
