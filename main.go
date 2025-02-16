package main

import (
	_ "github.com/Dimoonevs/hls-parser/pkg/parser/hls"
	"github.com/vharitonsky/iniflags"
	"hls-streamer/app/service"
	"log"
)

func main() {
	iniflags.Parse()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
