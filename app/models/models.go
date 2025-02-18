package models

type VideoFormat struct {
	URL        string `json:"url"`
	Resolution string `json:"size"`
	Bitrate    int    `json:"bitrate"`
}

type Video struct {
	ID       int            `json:"id"`
	Formats  []*VideoFormat `json:"formats"`
	Position int            `json:"position"`
}

type MediaPlaylist struct {
	Duration   float64
	URI        string
	Resolution string
}
