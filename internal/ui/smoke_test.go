package ui

import (
	"context"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func feed(m model, msg tea.Msg) model {
	nm, _ := m.Update(msg)
	return nm.(model)
}

func TestRenderSmoke(t *testing.T) {
	if testing.Short() {
		t.Skip("network smoke test; run without -short")
	}
	m := New()
	m = feed(m, tea.WindowSizeMsg{Width: 100, Height: 30})

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	games, err := m.client.Games(ctx, time.Now())
	if err != nil {
		t.Skipf("network unavailable: %v", err)
	}
	m = feed(m, gamesMsg{date: time.Now(), resp: games})
	out := m.View()
	if strings.TrimSpace(out) == "" {
		t.Fatal("empty scores view")
	}
	t.Log("\n=== SCORES ===\n" + out)

	if len(games.Leagues) == 0 {
		t.Skip("no leagues today")
	}
	lid := games.Leagues[0].ID
	ld, err := m.client.League(ctx, lid)
	if err != nil {
		t.Fatalf("league: %v", err)
	}
	m.leagueID = lid
	m = feed(m, leagueMsg{ld})
	m.screen = leagueScreen
	t.Log("\n=== STANDINGS ===\n" + m.View())
	m.leagueTab = 1
	t.Log("\n=== FIXTURES ===\n" + m.View())

	// game detail
	var gid string
	for _, lg := range games.Leagues {
		if len(lg.Games) > 0 {
			gid = lg.Games[0].ID
			break
		}
	}
	if gid == "" {
		return
	}
	gc, err := m.client.GameCenter(ctx, gid)
	if err != nil {
		t.Fatalf("gamecenter: %v", err)
	}
	m = feed(m, gameMsg{gc})
	m.screen = gameScreen
	t.Log("\n=== GAME ===\n" + m.View())
}
