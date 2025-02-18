package service

import (
	"flag"
	"fmt"
	"hls-streamer/app/repo/mysql"
	"hls-streamer/app/utils"
	"log"
	"net/http"
)

var (
	portStreamer = flag.Int("portStreamer", 8080, "The port streamer endpoint")
)

type Segment struct {
	Duration float64
	URI      string
}

func Run() error {
	ListVideo, err := mysql.GetConnection().GetVideoContent()
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
	for {
		for i := 0; i < len(ListVideo); i++ {
			stream.StartStream(&ListVideo[i])
		}
	}
}
