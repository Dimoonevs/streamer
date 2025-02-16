package hls

import (
	"fmt"
	"github.com/Dimoonevs/hls-parser/pkg/domain"
	"strings"
)

func GenerateMediaPlaylist(playlist *domain.MediaPlaylist) string {
	var builder strings.Builder

	builder.WriteString("#EXTM3U\n")
	builder.WriteString("#EXT-X-VERSION:3\n")
	builder.WriteString(fmt.Sprintf("#EXT-X-TARGETDURATION:%d\n", playlist.TargetDuration))
	builder.WriteString(fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%d\n", playlist.SeqNo))

	for _, segment := range playlist.Segments {
		builder.WriteString(fmt.Sprintf("#EXTINF:%.3f,\n%s\n", segment.Duration, segment.URI))
		if segment.Discontinuity {
			builder.WriteString("#EXT-X-DISCONTINUITY\n")
		}
	}

	return builder.String()
}
