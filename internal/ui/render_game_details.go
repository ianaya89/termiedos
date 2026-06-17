package ui

import (
	"strings"

	"termiedos/internal/api"

	"github.com/charmbracelet/lipgloss"
)

// renderGoals lists scorers grouped by team, e.g. "Brasil   Vinicius 32'".
func (m model) renderGoals(g api.GameDetail) string {
	hasAny := false
	for _, t := range g.Teams {
		if len(t.Goals) > 0 {
			hasAny = true
			break
		}
	}
	if !hasAny {
		return ""
	}

	nameW := 0
	for _, t := range g.Teams {
		if len(t.Goals) == 0 {
			continue
		}
		if w := lipgloss.Width(t.Name); w > nameW {
			nameW = w
		}
	}

	var lines []string
	for _, t := range g.Teams {
		if len(t.Goals) == 0 {
			continue
		}
		var parts []string
		for _, gl := range t.Goals {
			name := gl.PlayerSName
			if name == "" {
				name = gl.PlayerName
			}
			when := gl.TimeToDisplay
			parts = append(parts, styleGoal.Render(name)+" "+styleEventTime.Render(when))
		}
		line := " " + styleKVval.Render(padRight(t.Name, nameW)) + "  " + strings.Join(parts, ", ")
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// renderCards lists bookings chronologically, e.g. "▮ 37' Casemiro · Brasil".
func (m model) renderCards(g api.GameDetail) string {
	cards := g.Cards()
	if len(cards) == 0 {
		return ""
	}
	var lines []string
	for _, c := range cards {
		team := ""
		if c.Team == 1 && len(g.Teams) > 0 {
			team = g.Teams[0].Name
		} else if c.Team == 2 && len(g.Teams) > 1 {
			team = g.Teams[1].Name
		}
		player := ""
		if len(c.Texts) > 0 {
			player = c.Texts[0]
		}
		line := " " + cardMark(c.Type) + " " + styleEventTime.Render(padLeft(c.Time, 3)) +
			" " + styleKVval.Render(player)
		if team != "" {
			line += styleKV.Render(" · " + team)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// renderStats renders the match statistics as labelled proportional bars.
func (m model) renderStats(g api.GameDetail, width int) string {
	if len(g.Statistics) == 0 {
		return ""
	}
	nameW := 16
	valW := 4
	barW := width - nameW - 2*valW - 4
	if barW < 8 {
		barW = 8
	}
	if barW > 28 {
		barW = 28
	}

	var lines []string
	for _, s := range g.Statistics {
		if len(s.Values) < 2 {
			continue
		}
		home, away := s.Values[0], s.Values[1]
		p := 0.5
		if len(s.Percentages) >= 1 {
			p = s.Percentages[0]
		}
		line := " " + styleKV.Render(padRight(truncate(s.Name, nameW), nameW)) +
			styleKVval.Render(padLeft(home, valW)) + " " +
			statBar(p, barW) + " " +
			styleKVval.Render(padRight(away, valW))
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// statBar draws a width-cell bar split by the home share p (0..1).
func statBar(p float64, width int) string {
	if p < 0 {
		p = 0
	} else if p > 1 {
		p = 1
	}
	h := int(p*float64(width) + 0.5)
	if h > width {
		h = width
	}
	return styleBarHome.Render(strings.Repeat("█", h)) +
		styleBarAway.Render(strings.Repeat("░", width-h))
}
