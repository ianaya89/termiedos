package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	baseURL   = "https://api.promiedos.com.ar"
	clientVer = "1.11.7.5"
	userAgent = "Mozilla/5.0 (promiedos-tui)"
)

type Client struct {
	http *http.Client
}

func New() *Client {
	return &Client{http: &http.Client{Timeout: 15 * time.Second}}
}

func (c *Client) get(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-VER", clientVer)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "https://www.promiedos.com.ar/")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %s", path, resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// DateFmt is the path date layout used by /games/{date}.
const DateFmt = "02-01-2006"

func (c *Client) Games(ctx context.Context, day time.Time) (*GamesResponse, error) {
	var r GamesResponse
	if err := c.get(ctx, "/games/"+day.Format(DateFmt), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Client) League(ctx context.Context, id string) (*LeagueData, error) {
	var r LeagueData
	if err := c.get(ctx, "/league/tables_and_fixtures/"+id, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Client) GameCenter(ctx context.Context, id string) (*GameCenter, error) {
	var r GameCenter
	if err := c.get(ctx, "/gamecenter/"+id, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
