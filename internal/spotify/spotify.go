package spotify

import (
	"context"
	"fmt"

	"github.com/Artem/DC_bot/internal/player"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

type Client struct {
	client *spotify.Client
}

type SearchResult struct {
	Title    string
	Artist   string
	URL      string
	Duration string
	ID       spotify.ID
}

func New(clientID, clientSecret string) (*Client, error) {
	ctx := context.Background()

	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить токен: %w", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	return &Client{client: client}, nil
}

func (c *Client) Search(query string) ([]SearchResult, error) {
	ctx := context.Background()

	results, err := c.client.Search(ctx, query, spotify.SearchTypeTrack, spotify.Limit(10))
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска: %w", err)
	}

	var tracks []SearchResult
	if results.Tracks != nil {
		for _, track := range results.Tracks.Tracks {
			artists := ""
			for i, artist := range track.Artists {
				if i > 0 {
					artists += ", "
				}
				artists += artist.Name
			}

			duration := fmt.Sprintf("%d:%02d", track.Duration/60000, (track.Duration/1000)%60)

			tracks = append(tracks, SearchResult{
				Title:    track.Name,
				Artist:   artists,
				URL:      track.ExternalURLs["spotify"],
				Duration: duration,
				ID:       track.ID,
			})
		}
	}

	return tracks, nil
}

func (c *Client) GetTrack(id spotify.ID) (*player.Track, error) {
	ctx := context.Background()

	track, err := c.client.GetTrack(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения трека: %w", err)
	}

	artists := ""
	for i, artist := range track.Artists {
		if i > 0 {
			artists += ", "
		}
		artists += artist.Name
	}

	duration := fmt.Sprintf("%d:%02d", track.Duration/60000, (track.Duration/1000)%60)

	return &player.Track{
		Title:    fmt.Sprintf("%s - %s", track.Name, artists),
		URL:      track.ExternalURLs["spotify"],
		Platform: "spotify",
		Duration: duration,
	}, nil
}
