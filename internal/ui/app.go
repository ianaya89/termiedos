package ui

import (
	"context"
	"time"

	"termiedos/internal/api"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen int

const (
	scoresScreen screen = iota
	leagueScreen
	gameScreen
)

type focusArea int

const (
	focusMain focusArea = iota
	focusSidebar
)

const sidebarWidth = 24

type scoreRow struct {
	header bool
	league *api.League
	game   *api.Game
}

type model struct {
	client  *api.Client
	w, h    int
	screen  screen
	loading bool
	err     error
	spinner spinner.Model

	// scores
	date    time.Time
	games   *api.GamesResponse
	rows    []scoreRow
	cursor  int
	offset  int
	focus   focusArea
	sideIdx int

	// league
	leagueID  string
	league    *api.LeagueData
	leagueTab int // 0 = standings, 1 = fixtures
	fixRound  int
	lgCursor  int
	lgOffset  int

	// game
	game       *api.GameCenter
	gameFrom   screen
	gameOffset int
}

func New() model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(colAccent)
	return model{
		client:  api.New(),
		date:    time.Now(),
		screen:  scoresScreen,
		loading: true,
		spinner: sp,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchGames(m.date), tickCmd())
}

// ---- messages ----

type gamesMsg struct {
	date time.Time
	resp *api.GamesResponse
}
type leagueMsg struct{ resp *api.LeagueData }
type gameMsg struct{ resp *api.GameCenter }
type errMsg struct{ err error }
type tickMsg struct{}

func tickCmd() tea.Cmd {
	return tea.Tick(15*time.Second, func(time.Time) tea.Msg { return tickMsg{} })
}

func (m model) fetchGames(day time.Time) tea.Cmd {
	c := m.client
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		r, err := c.Games(ctx, day)
		if err != nil {
			return errMsg{err}
		}
		return gamesMsg{date: day, resp: r}
	}
}

func (m model) fetchLeague(id string) tea.Cmd {
	c := m.client
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		r, err := c.League(ctx, id)
		if err != nil {
			return errMsg{err}
		}
		return leagueMsg{r}
	}
}

func (m model) fetchGame(id string) tea.Cmd {
	c := m.client
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		r, err := c.GameCenter(ctx, id)
		if err != nil {
			return errMsg{err}
		}
		return gameMsg{r}
	}
}

// refreshCurrent silently refetches the data backing the active screen.
func (m model) refreshCurrent() tea.Cmd {
	switch m.screen {
	case leagueScreen:
		if m.leagueID != "" {
			return m.fetchLeague(m.leagueID)
		}
	case gameScreen:
		if m.game != nil {
			return m.fetchGame(m.game.Game.ID)
		}
	default:
		return m.fetchGames(m.date)
	}
	return nil
}

func (m *model) buildRows() {
	m.rows = m.rows[:0]
	if m.games == nil {
		return
	}
	for i := range m.games.Leagues {
		lg := &m.games.Leagues[i]
		m.rows = append(m.rows, scoreRow{header: true, league: lg})
		for j := range lg.Games {
			m.rows = append(m.rows, scoreRow{league: lg, game: &lg.Games[j]})
		}
	}
	if m.cursor >= len(m.rows) {
		m.cursor = len(m.rows) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tickMsg:
		return m, tea.Batch(m.refreshCurrent(), tickCmd())

	case errMsg:
		m.loading = false
		m.err = msg.err
		return m, nil

	case gamesMsg:
		m.loading = false
		m.err = nil
		m.date = msg.date
		m.games = msg.resp
		m.buildRows()
		if m.sideIdx >= len(m.games.Leagues) {
			m.sideIdx = 0
		}
		return m, nil

	case leagueMsg:
		m.loading = false
		m.err = nil
		m.league = msg.resp
		m.fixRound = selectedRound(msg.resp)
		return m, nil

	case gameMsg:
		m.loading = false
		m.err = nil
		m.game = msg.resp
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "r":
		m.loading = true
		return m, m.refreshCurrent()
	}
	switch m.screen {
	case scoresScreen:
		return m.keyScores(msg)
	case leagueScreen:
		return m.keyLeague(msg)
	case gameScreen:
		return m.keyGame(msg)
	}
	return m, nil
}

func (m model) keyScores(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "left", "h":
		m.date = m.date.AddDate(0, 0, -1)
		m.loading, m.cursor, m.offset = true, 0, 0
		return m, m.fetchGames(m.date)
	case "right", "l":
		m.date = m.date.AddDate(0, 0, 1)
		m.loading, m.cursor, m.offset = true, 0, 0
		return m, m.fetchGames(m.date)
	case "t":
		m.date = time.Now()
		m.loading, m.cursor, m.offset = true, 0, 0
		return m, m.fetchGames(m.date)
	case "tab":
		if m.focus == focusMain {
			m.focus = focusSidebar
		} else {
			m.focus = focusMain
		}
		return m, nil
	case "up", "k":
		if m.focus == focusSidebar {
			if m.sideIdx > 0 {
				m.sideIdx--
			}
		} else {
			m.moveCursor(-1)
		}
		return m, nil
	case "down", "j":
		if m.focus == focusSidebar {
			if m.games != nil && m.sideIdx < len(m.games.Leagues)-1 {
				m.sideIdx++
			}
		} else {
			m.moveCursor(1)
		}
		return m, nil
	case "enter":
		if m.focus == focusSidebar {
			if m.games != nil && m.sideIdx < len(m.games.Leagues) {
				return m.openLeague(m.games.Leagues[m.sideIdx].ID)
			}
			return m, nil
		}
		if m.cursor < len(m.rows) {
			row := m.rows[m.cursor]
			if row.header {
				return m.openLeague(row.league.ID)
			}
			return m.openGame(row.game.ID, scoresScreen)
		}
	}
	return m, nil
}

func (m *model) moveCursor(d int) {
	n := len(m.rows)
	if n == 0 {
		return
	}
	m.cursor += d
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= n {
		m.cursor = n - 1
	}
}

func (m model) openLeague(id string) (model, tea.Cmd) {
	m.screen = leagueScreen
	m.leagueID = id
	m.league = nil
	m.leagueTab = 0
	m.lgCursor, m.lgOffset = 0, 0
	m.loading = true
	return m, m.fetchLeague(id)
}

func (m model) openGame(id string, from screen) (model, tea.Cmd) {
	m.screen = gameScreen
	m.gameFrom = from
	m.game = nil
	m.gameOffset = 0
	m.loading = true
	return m, m.fetchGame(id)
}

func (m model) keyLeague(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "backspace", "b":
		m.screen = scoresScreen
		m.loading = false
		return m, nil
	case "tab", "1", "2":
		if msg.String() == "1" {
			m.leagueTab = 0
		} else if msg.String() == "2" {
			m.leagueTab = 1
		} else {
			m.leagueTab = (m.leagueTab + 1) % 2
		}
		m.lgCursor, m.lgOffset = 0, 0
		return m, nil
	case "up", "k":
		if m.lgCursor > 0 {
			m.lgCursor--
		}
		return m, nil
	case "down", "j":
		m.lgCursor++ // clamped at render
		return m, nil
	case "left", "h":
		if m.leagueTab == 1 {
			m.prevRound()
			m.lgCursor, m.lgOffset = 0, 0
		}
		return m, nil
	case "right", "l":
		if m.leagueTab == 1 {
			m.nextRound()
			m.lgCursor, m.lgOffset = 0, 0
		}
		return m, nil
	case "enter":
		if m.leagueTab == 1 && m.league != nil {
			gs := roundGames(m.league, m.fixRound)
			if m.lgCursor < len(gs) {
				return m.openGame(gs[m.lgCursor].ID, leagueScreen)
			}
		}
		return m, nil
	}
	return m, nil
}

func (m *model) prevRound() {
	if m.league == nil {
		return
	}
	for i := m.fixRound - 1; i >= 0; i-- {
		if len(m.league.Games.Filters[i].Games) > 0 {
			m.fixRound = i
			return
		}
	}
}

func (m *model) nextRound() {
	if m.league == nil {
		return
	}
	for i := m.fixRound + 1; i < len(m.league.Games.Filters); i++ {
		if len(m.league.Games.Filters[i].Games) > 0 {
			m.fixRound = i
			return
		}
	}
}

func (m model) keyGame(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	max := m.gameMaxOffset(m.w, m.gameViewHeight())
	switch msg.String() {
	case "esc", "backspace", "b":
		m.screen = m.gameFrom
		m.loading = false
		return m, nil
	case "up", "k":
		if m.gameOffset > 0 {
			m.gameOffset--
		}
	case "down", "j":
		if m.gameOffset < max {
			m.gameOffset++
		}
	case "pgup":
		m.gameOffset -= 10
	case "pgdown", " ":
		m.gameOffset += 10
	case "g", "home":
		m.gameOffset = 0
	case "G", "end":
		m.gameOffset = max
	}
	if m.gameOffset < 0 {
		m.gameOffset = 0
	}
	if m.gameOffset > max {
		m.gameOffset = max
	}
	return m, nil
}

// gameViewHeight mirrors View()'s body height (full height minus header + help).
func (m model) gameViewHeight() int {
	h := m.h - 2
	if h < 1 {
		h = 1
	}
	return h
}
