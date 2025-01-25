package service

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var playlistContent string
var mutex sync.RWMutex
var sleepTime float64

type Segment struct {
	Duration float64
	URI      string
}

func Run() {
	playlistURL := "http://your-video-service.pp.ua/video/service/9552a004606fe3cff7f21fad75703318/company-name_2_1737821675_live4m.mp40.m3u8"

	segments, targetDuration, err := getPlaylistData(playlistURL)
	if err != nil {
		fmt.Println("Ошибка получения данных плейлиста:", err)
		return
	}

	fmt.Printf("Target duration: %.1f seconds\n", targetDuration)
	for _, segment := range segments {
		fmt.Printf("Segment: URI=%s, Duration=%.1f\n", segment.URI, segment.Duration)
	}

	go generatePlaylist(segments, targetDuration)
	go func() {
		for {
			log.Println(playlistContent)
			time.Sleep(time.Duration(sleepTime+1) * time.Second)
		}
	}()

	http.HandleFunc("/playlist.m3u8", servePlaylist)
	fmt.Println("Сервер запущен на http://localhost:8080/playlist.m3u8")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}

func getPlaylistData(playlistURL string) ([]Segment, float64, error) {
	resp, err := http.Get(playlistURL)
	if err != nil {
		return nil, 0, fmt.Errorf("не удалось получить плейлист: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("ошибка загрузки плейлиста: статус %d", resp.StatusCode)
	}

	var segments []Segment
	var targetDuration float64
	var lastDuration float64

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#EXT-X-TARGETDURATION:") {
			durationStr := strings.TrimPrefix(line, "#EXT-X-TARGETDURATION:")
			targetDuration, err = strconv.ParseFloat(durationStr, 64)
			if err != nil {
				return nil, 0, fmt.Errorf("ошибка парсинга TARGETDURATION: %v", err)
			}
		} else if strings.HasPrefix(line, "#EXTINF:") {
			durationStr := strings.TrimSuffix(strings.TrimPrefix(line, "#EXTINF:"), ",")
			lastDuration, err = strconv.ParseFloat(durationStr, 64)
			if err != nil {
				return nil, 0, fmt.Errorf("ошибка парсинга EXTINF: %v", err)
			}
		} else if strings.HasSuffix(line, ".ts") {
			segments = append(segments, Segment{
				Duration: lastDuration,
				URI:      line,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, fmt.Errorf("ошибка чтения плейлиста: %v", err)
	}

	return segments, targetDuration, nil
}

func generatePlaylist(segments []Segment, targetDuration float64) {
	if len(segments) < 5 {
		fmt.Println("Недостаточно сегментов в плейлисте")
		return
	}

	startIndex := 0
	sequenceNumber := 0
	lastDurations := []float64{}

	for {
		playlist := "#EXTM3U\n#EXT-X-VERSION:3\n"
		playlist += fmt.Sprintf("#EXT-X-TARGETDURATION:%.1f\n", targetDuration)
		playlist += "#EXT-X-MEDIA-SEQUENCE:" + fmt.Sprintf("%d", sequenceNumber) + "\n"

		for i := 0; i < 5; i++ {
			index := (startIndex + i) % len(segments)
			playlist += fmt.Sprintf("#EXTINF:%.1f,\n%s\n", segments[index].Duration, segments[index].URI)

			if len(lastDurations) == 3 {
				lastDurations = lastDurations[1:]
			}
			lastDurations = append(lastDurations, segments[index].Duration)
		}

		mutex.Lock()
		playlistContent = playlist
		mutex.Unlock()

		sleepDuration := calculateAverage(lastDurations)
		fmt.Printf("Средняя длительность: %.2f секунд. Засыпаем...\n", sleepDuration)
		sleepTime = sleepDuration
		sequenceNumber++
		startIndex = (startIndex + 1) % len(segments)

		time.Sleep(time.Duration(sleepDuration) * time.Second)
	}
}

func calculateAverage(durations []float64) float64 {
	if len(durations) == 0 {
		return 0
	}
	var sum float64
	for _, d := range durations {
		sum += d
	}
	return sum / float64(len(durations))
}

func servePlaylist(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(playlistContent))
}
