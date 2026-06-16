package ui

import (
	"strings"

	"termiedos/internal/api"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderScores(width, height int) string {
	if m.err != nil {
		return styleErr.Width(width).Render("Error: " + m.err.Error())
	}
	if m.games == nil {
		return lipgloss.NewStyle().Width(width).Padding(1, 2).Foreground(colMuted).Render(m.spinner.View() + " cargando…")
	}
	if len(m.rows) == 0 {
		return lipgloss.NewStyle().Width(width).Padding(1, 2).Foreground(colMuted).Render("Sin partidos para esta fecha.")
	}

	start := 0
	if m.cursor >= height {
		start = m.cursor - height + 1
	}
	end := start + height
	if end > len(m.rows) {
		end = len(m.rows)
	}

	var b strings.Builder
	for i := start; i < end; i++ {
		row := m.rows[i]
		sel := m.focus == focusMain && i == m.cursor
		if row.header {
			b.WriteString(m.renderLeagueBar(*row.league, width))
		} else {
			b.WriteString(m.renderGameRow(*row.game, width, sel))
		}
		b.WriteByte('\n')
	}
	return lipgloss.NewStyle().Width(width).Render(b.String())
}

func (m model) renderLeagueBar(lg api.League, width int) string {
	name := lg.Name
	country := ""
	if lg.CountryName != "" && lg.CountryName != lg.Name {
		country = lg.CountryName
	}
	label := "▌ " + name
	if country != "" {
		label += "  " + styleLeagueCountry.Render(country)
	}
	return styleLeagueHeader.Width(width).Render(label)
}

func (m model) renderGameRow(g api.Game, width int, sel bool) string {
	if len(g.Teams) < 2 {
		return ""
	}
	home, away := g.Teams[0], g.Teams[1]

	clockW := 8
	scoreW := 7
	nameW := (width - clockW - scoreW - 4) / 2
	if nameW < 6 {
		nameW = 6
	}

	clock := g.Clock()
	var clockS string
	switch {
	case g.IsLive():
		clockS = styleClockLive.Render(padRight(clock, clockW))
	case g.IsFinal():
		clockS = styleClockFinal.Render(padRight(clock, clockW))
	default:
		clockS = styleClock.Render(padRight(clock, clockW))
	}

	hStyle, aStyle := styleTeam, styleTeam
	if g.IsFinal() || g.HasScore() {
		switch g.Winner {
		case 0:
			hStyle, aStyle = styleTeamWin, styleTeamLose
		case 1:
			hStyle, aStyle = styleTeamLose, styleTeamWin
		}
	}

	hName := hStyle.Render(padLeft(home.Name, nameW))
	aName := aStyle.Render(padRight(away.Name, nameW))
	hBlock := teamColorBlock(home.Colors.Color)
	aBlock := teamColorBlock(away.Colors.Color)

	score := g.ScoreText()
	scoreStyle := styleScore
	if g.IsLive() {
		scoreStyle = styleScoreLive
	} else if !g.HasScore() {
		scoreStyle = styleClock
	}
	scoreS := scoreStyle.Render(center(score, scoreW))

	line := " " + clockS + hName + " " + hBlock + scoreS + aBlock + " " + aName
	if sel {
		return styleRowSel.Width(width).Render(line)
	}
	return lipgloss.NewStyle().Width(width).Render(line)
}
