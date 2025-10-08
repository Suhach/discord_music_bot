package youtube

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"discord_music_bot/internal/player"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Client struct {
	service *youtube.Service
}

type SearchResult struct {
	Title       string
	ChannelName string
	VideoID     string
	Duration    string
	URL         string
}

func New(apiKey string) (*Client, error) {
	ctx := context.Background()

	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать YouTube клиент: %w", err)
	}

	return &Client{service: service}, nil
}

func (c *Client) Search(query string) ([]SearchResult, error) {
	call := c.service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(10).
		Type("video")

	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска: %w", err)
	}

	var results []SearchResult
	var videoIDs []string

	for _, item := range response.Items {
		videoIDs = append(videoIDs, item.Id.VideoId)
	}

	// Получаем информацию о длительности видео
	videosCall := c.service.Videos.List([]string{"contentDetails", "snippet"}).
		Id(videoIDs...)

	videosResponse, err := videosCall.Do()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения деталей видео: %w", err)
	}

	for _, video := range videosResponse.Items {
		duration := parseDuration(video.ContentDetails.Duration)

		results = append(results, SearchResult{
			Title:       video.Snippet.Title,
			ChannelName: video.Snippet.ChannelTitle,
			VideoID:     video.Id,
			Duration:    duration,
			URL:         fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.Id),
		})
	}

	return results, nil
}

func (c *Client) GetTrack(videoID string) (*player.Track, error) {
	call := c.service.Videos.List([]string{"snippet", "contentDetails"}).
		Id(videoID)

	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения видео: %w", err)
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("видео не найдено")
	}

	video := response.Items[0]
	duration := parseDuration(video.ContentDetails.Duration)

	return &player.Track{
		Title:    video.Snippet.Title,
		URL:      fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
		Platform: "youtube",
		Duration: duration,
	}, nil
}

// parseDuration парсит ISO 8601 длительность в читаемый формат
func parseDuration(duration string) string {
	re := regexp.MustCompile(`PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`)
	matches := re.FindStringSubmatch(duration)

	if len(matches) == 0 {
		return "0:00"
	}

	hours := 0
	minutes := 0
	seconds := 0

	if matches[1] != "" {
		hours, _ = strconv.Atoi(matches[1])
	}
	if matches[2] != "" {
		minutes, _ = strconv.Atoi(matches[2])
	}
	if matches[3] != "" {
		seconds, _ = strconv.Atoi(matches[3])
	}

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
