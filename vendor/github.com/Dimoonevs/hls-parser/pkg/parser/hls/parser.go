package hls

import (
	"bufio"
	"fmt"
	"github.com/Dimoonevs/hls-parser/pkg/domain"
	"net/http"
	"strconv"
	"strings"
)

func ParseMediaPlaylist(url string) (*domain.MediaPlaylist, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки плейлиста: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("сервер вернул код %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	return parsePlaylist(scanner)
}

func parsePlaylist(scanner *bufio.Scanner) (playlist *domain.MediaPlaylist, err error) {
	playlist = &domain.MediaPlaylist{
		Segments: []domain.Segment{},
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "#EXT-X-TARGETDURATION:") {
			value, err := strconv.Atoi(strings.TrimPrefix(line, "#EXT-X-TARGETDURATION:"))
			if err != nil {
				return nil, fmt.Errorf("ошибка парсинга TARGETDURATION: %v", err)
			}
			playlist.TargetDuration = value
		} else if strings.HasPrefix(line, "#EXTINF:") {
			segment, err := parseSegment(line, scanner)
			if err != nil {
				return nil, err
			}
			playlist.Segments = append(playlist.Segments, *segment)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %v", err)
	}

	if len(playlist.Segments) > 0 {
		playlist.Segments[len(playlist.Segments)-1].Discontinuity = true
	}

	return playlist, nil
}

func parseSegment(extinfLine string, scanner *bufio.Scanner) (*domain.Segment, error) {
	parts := strings.Split(extinfLine, ":")
	if len(parts) < 2 {
		return nil, fmt.Errorf("неверный формат EXTINF: %s", extinfLine)
	}

	durationStr := strings.TrimSuffix(parts[1], ",")
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга EXTINF длительности: %v", err)
	}

	if !scanner.Scan() {
		return nil, fmt.Errorf("отсутствует сегментный файл после EXTINF")
	}
	segmentURI := strings.TrimSpace(scanner.Text())

	return &domain.Segment{
		Duration: duration,
		URI:      segmentURI,
	}, nil
}
