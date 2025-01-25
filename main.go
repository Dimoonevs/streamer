package main

import (
	"github.com/vharitonsky/iniflags"
	"hls-streamer/app/service"
)

func main() {
	iniflags.Parse()
	service.Run()
}
