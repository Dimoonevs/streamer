package domain

type MediaPlaylist struct {
	TargetDuration int
	Segments       []Segment
	Size           string
	SeqNo          int
}

type Segment struct {
	Duration      float64
	URI           string
	Discontinuity bool
}
type MasterPlaylist struct {
}

type MediaPlaylistForCreate struct {
	MediaURL string
	Size     string
}
