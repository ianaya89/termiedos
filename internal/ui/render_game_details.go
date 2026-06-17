package ui

import (
	"strconv"
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

// renderForm shows each team's recent results as G/E/P chips.
func (m model) renderForm(g api.GameDetail) string {
	rf := g.RecentForm
	if len(rf.Home) == 0 && len(rf.Away) == 0 {
		return ""
	}
	if len(g.Teams) < 2 {
		return ""
	}
	nameW := lipgloss.Width(g.Teams[0].Name)
	if w := lipgloss.Width(g.Teams[1].Name); w > nameW {
		nameW = w
	}
	row := func(name string, codes []int) string {
		var chips []string
		for _, c := range codes {
			chips = append(chips, formMark(c))
		}
		return " " + styleKVval.Render(padRight(name, nameW)) + "  " + strings.Join(chips, " ")
	}
	return row(g.Teams[0].Name, rf.Home) + "\n" + row(g.Teams[1].Name, rf.Away)
}

// renderLineups shows both starting XIs side by side with event markers.
func (m model) renderLineups(g api.GameDetail, width int) string {
	teams := g.Players.Lineups.Teams
	if len(teams) < 2 || len(g.Teams) < 2 {
		return ""
	}
	colW := (width - 2) / 2
	if colW < 16 {
		colW = 16
	}

	col := func(team api.LineupTeam, name string) string {
		head := name
		if team.Formation != "" {
			head += " (" + team.Formation + ")"
		}
		lines := []string{styleKVval.Bold(true).Render(truncate(head, colW))}
		if c := team.Coach(); c != "" {
			lines = append(lines, styleKV.Render(padRight("DT: "+truncate(c, colW-4), colW)))
		}
		for _, p := range team.Starting {
			lines = append(lines, playerLine(p, colW))
		}
		return strings.Join(lines, "\n")
	}

	left := col(teams[0], g.Teams[0].Name)
	right := col(teams[1], g.Teams[1].Name)
	return lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(colW).Render(left),
		"  ",
		lipgloss.NewStyle().Width(colW).Render(right))
}

func playerLine(p api.LineupPlayer, width int) string {
	name := p.ShortName
	if name == "" {
		name = p.Name
	}
	var marks string
	if p.Events.Goals.Goals > 0 || p.Events.Goals.OwnGoals > 0 {
		marks += "⚽"
	}
	if p.Events.Cards.Red {
		marks += redMark
	} else if p.Events.Cards.Yellow {
		marks += yellowMark
	}
	if p.Events.Substitution.HasSubstitution {
		marks += styleEventTime.Render("↓")
	}
	cap := ""
	if p.IsCaptain {
		cap = styleCaptain.Render(" (C)")
	}
	num := styleNum.Render(padLeft(itoa(p.JerseyNum), 2))
	avail := width - 3 - lipgloss.Width(marks) - lipgloss.Width(cap)
	if avail < 4 {
		avail = 4
	}
	return " " + num + " " + styleGoal.Render(truncate(name, avail)) + cap + " " + marks
}

// renderH2H shows the head-to-head balance and recent meetings.
func (m model) renderH2H(g api.GameDetail) string {
	h := g.HeadToHead
	if h.HomeWins == 0 && h.AwayWins == 0 && h.Draws == 0 && len(h.Games) == 0 {
		return ""
	}
	if len(g.Teams) < 2 {
		return ""
	}
	summary := " " + styleKVval.Render(g.Teams[0].Name) + styleKV.Render(" "+itoa(h.HomeWins)) +
		styleKV.Render(" · "+itoa(h.Draws)+" E · ") +
		styleKV.Render(itoa(h.AwayWins)+" ") + styleKVval.Render(g.Teams[1].Name)
	lines := []string{summary}
	shown := 0
	for _, pg := range h.Games {
		if pg.ID == g.ID || len(pg.Teams) < 2 || !pg.HasScore() {
			continue
		}
		date := ""
		if len(pg.StartTime) >= 10 {
			date = pg.StartTime[:10]
		}
		line := " " + styleEventTime.Render(date) + "  " +
			styleKV.Render(truncate(pg.Teams[0].Name, 14)) + " " +
			styleKVval.Render(pg.ScoreText()) + " " +
			styleKV.Render(truncate(pg.Teams[1].Name, 14))
		lines = append(lines, line)
		if shown++; shown >= 5 {
			break
		}
	}
	return strings.Join(lines, "\n")
}

func itoa(n int) string { return strconv.Itoa(n) }

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
