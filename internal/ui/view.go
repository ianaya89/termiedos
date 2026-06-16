package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = model{}

func (m model) View() string {
	if m.w == 0 {
		return "loading…"
	}
	header := m.renderHeader()
	help := m.renderHelp()
	bodyH := m.h - lipgloss.Height(header) - lipgloss.Height(help)
	if bodyH < 1 {
		bodyH = 1
	}

	var body string
	switch m.screen {
	case leagueScreen:
		body = m.renderLeague(m.w, bodyH)
	case gameScreen:
		body = m.renderGame(m.w, bodyH)
	default:
		side := m.renderSidebar(bodyH)
		main := m.renderScores(m.w-sidebarWidth, bodyH)
		body = lipgloss.JoinHorizontal(lipgloss.Top, side, main)
	}
	body = lipgloss.NewStyle().Height(bodyH).MaxHeight(bodyH).Render(body)
	return lipgloss.JoinVertical(lipgloss.Left, header, body, help)
}

func (m model) renderHeader() string {
	left := styleTitle.Render("⚽ TERMIEDOS")
	var mid string
	switch m.screen {
	case leagueScreen:
		name := "Liga"
		if m.league != nil {
			name = m.league.League.Name
		}
		mid = styleTitleDate.Render(name)
	case gameScreen:
		t := "Partido"
		if m.game != nil && len(m.game.Game.Teams) == 2 {
			t = m.game.Game.Teams[0].Name + " vs " + m.game.Game.Teams[1].Name
		}
		mid = styleTitleDate.Render(t)
	default:
		label := m.date.Format("Mon 02-01-2006")
		mid = styleTitleDate.Render("◀ " + label + " ▶")
	}

	right := ""
	if m.loading {
		right = m.spinner.View()
	} else if n := m.liveCount(); n > 0 {
		right = styleClockLive.Render(fmt.Sprintf("● %d EN VIVO", n))
	}

	gap := m.w - lipgloss.Width(left) - lipgloss.Width(mid) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}
	bar := left + mid + strings.Repeat(" ", gap) + right
	return lipgloss.NewStyle().Background(colHeader).MaxWidth(m.w).Render(bar)
}

func (m model) liveCount() int {
	n := 0
	if m.games == nil {
		return 0
	}
	for _, lg := range m.games.Leagues {
		for _, g := range lg.Games {
			if g.IsLive() {
				n++
			}
		}
	}
	return n
}

func (m model) renderHelp() string {
	var keys string
	switch m.screen {
	case leagueScreen:
		keys = "tab posiciones/fixture · ←→ fecha · ↑↓ scroll · enter partido · b atrás · r recargar · q salir"
	case gameScreen:
		keys = "b atrás · r recargar · q salir"
	default:
		keys = "↑↓ mover · ←→ día · t hoy · tab panel · enter abrir · r recargar · q salir"
	}
	return styleHelp.Width(m.w).Render(keys)
}

func (m model) renderSidebar(height int) string {
	var b strings.Builder
	b.WriteString(styleSidebarTitle.Render("LIGAS") + "\n")
	rows := height - 1
	if m.games != nil {
		for i, lg := range m.games.Leagues {
			if i >= rows {
				break
			}
			name := truncate(lg.Name, sidebarWidth-3)
			if m.focus == focusSidebar && i == m.sideIdx {
				b.WriteString(styleSidebarSel.Width(sidebarWidth).Render("› "+name) + "\n")
			} else {
				b.WriteString(styleSidebarItem.Width(sidebarWidth).Render("  "+name) + "\n")
			}
		}
	}
	panel := lipgloss.NewStyle().
		Width(sidebarWidth).Height(height).
		Background(colPanel).
		BorderStyle(lipgloss.NormalBorder()).BorderRight(true).BorderForeground(colBorder)
	return panel.Render(b.String())
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	if n <= 1 {
		return string(r[:n])
	}
	return string(r[:n-1]) + "…"
}

func padRight(s string, n int) string {
	w := lipgloss.Width(s)
	if w >= n {
		return truncate(s, n)
	}
	return s + strings.Repeat(" ", n-w)
}

func padLeft(s string, n int) string {
	w := lipgloss.Width(s)
	if w >= n {
		return truncate(s, n)
	}
	return strings.Repeat(" ", n-w) + s
}

func center(s string, n int) string {
	w := lipgloss.Width(s)
	if w >= n {
		return s
	}
	l := (n - w) / 2
	return strings.Repeat(" ", l) + s + strings.Repeat(" ", n-w-l)
}
