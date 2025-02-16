package hls

import (
	"fmt"
	"github.com/Dimoonevs/hls-parser/pkg/domain"
	"strings"
)

func GenerateMasterPlaylist(playlists []domain.MediaPlaylistForCreate) string {
	var builder strings.Builder

	builder.WriteString("#EXTM3U\n")
	builder.WriteString("#EXT-X-VERSION:3\n")

	for _, playlist := range playlists {

		bitrate := estimateBitrate(playlist.Size)

		builder.WriteString(fmt.Sprintf(
			"#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%s\n%s\n",
			bitrate, playlist.Size, playlist.MediaURL,
		))
	}

	return builder.String()
}

func estimateBitrate(size string) int {
	switch size {
	case "1920x1080":
		return 5500
	case "1280x720":
		return 3400
	case "854x480":
		return 1200
	default:
		return 1000000
	}
}
