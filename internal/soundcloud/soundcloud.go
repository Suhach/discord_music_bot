package soundcloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"discord_music_bot/internal/player"
)

type Client struct {
	clientID string
}

type SearchResult struct {
	Title    string
	Artist   string
	URL      string
	Duration string
	ID       int64
}

type scTrack struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	PermalinkURL string `json:"permalink_url"`
	Duration     int64  `json:"duration"`
	User         struct {
		Username string `json:"username"`
	} `json:"user"`
}

type scSearchResponse struct {
	Collection []scTrack `json:"collection"`
}

func New(clientID string) *Client {
	return &Client{clientID: clientID}
}

func (c *Client) Search(query string) ([]SearchResult, error) {
	encodedQuery := url.QueryEscape(query)
	apiURL := fmt.Sprintf("https://api-v2.soundcloud.com/search/tracks?q=%s&client_id=%s&limit=10", encodedQuery, c.clientID)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка API: статус %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var searchResp scSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	var results []SearchResult
	for _, track := range searchResp.Collection {
		duration := formatDuration(track.Duration)

		results = append(results, SearchResult{
			Title:    track.Title,
			Artist:   track.User.Username,
			URL:      track.PermalinkURL,
			Duration: duration,
			ID:       track.ID,
		})
	}

	return results, nil
}

func (c *Client) GetTrack(trackID int64) (*player.Track, error) {
	apiURL := fmt.Sprintf("https://api-v2.soundcloud.com/tracks/%d?client_id=%s", trackID, c.clientID)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка API: статус %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var track scTrack
	if err := json.Unmarshal(body, &track); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	duration := formatDuration(track.Duration)

	return &player.Track{
		Title:    fmt.Sprintf("%s - %s", track.Title, track.User.Username),
		URL:      track.PermalinkURL,
		Platform: "soundcloud",
		Duration: duration,
	}, nil
}

func formatDuration(milliseconds int64) string {
	seconds := milliseconds / 1000
	minutes := seconds / 60
	seconds = seconds % 60

	if minutes >= 60 {
		hours := minutes / 60
		minutes = minutes % 60
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}

	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
