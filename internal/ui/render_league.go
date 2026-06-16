package ui

import (
	"strconv"
	"strings"

	"termiedos/internal/api"

	"github.com/charmbracelet/lipgloss"
)

func selectedRound(ld *api.LeagueData) int {
	if ld == nil {
		return 0
	}
	for i, f := range ld.Games.Filters {
		if f.Selected && len(f.Games) > 0 {
			return i
		}
	}
	for i, f := range ld.Games.Filters {
		if len(f.Games) > 0 {
			return i
		}
	}
	return 0
}

func roundGames(ld *api.LeagueData, idx int) []api.Game {
	if ld == nil || idx < 0 || idx >= len(ld.Games.Filters) {
		return nil
	}
	return ld.Games.Filters[idx].Games
}

func (m model) renderLeague(width, height int) string {
	if m.err != nil {
		return styleErr.Width(width).Render("Error: " + m.err.Error())
	}
	if m.league == nil {
		return lipgloss.NewStyle().Width(width).Padding(1, 2).Foreground(colMuted).Render(m.spinner.View() + " cargando…")
	}

	tabs := m.renderTabs(width)
	bodyH := height - lipgloss.Height(tabs)
	if bodyH < 1 {
		bodyH = 1
	}
	var body string
	if m.leagueTab == 0 {
		body = m.renderStandings(width, bodyH)
	} else {
		body = m.renderFixtures(width, bodyH)
	}
	return lipgloss.JoinVertical(lipgloss.Left, tabs, body)
}

func (m model) renderTabs(width int) string {
	var pos, fix string
	if m.leagueTab == 0 {
		pos = styleTabActive.Render("Posiciones")
		fix = styleTabInactive.Render("Fixture")
	} else {
		pos = styleTabInactive.Render("Posiciones")
		fix = styleTabActive.Render("Fixture")
	}
	bar := pos + " " + fix
	return lipgloss.NewStyle().Width(width).Background(colPanel).Render(bar)
}

func (m model) renderStandings(width, height int) string {
	var lines []string
	if len(m.league.TablesGroups) == 0 {
		return lipgloss.NewStyle().Width(width).Padding(1, 2).Foreground(colMuted).Render("Sin tabla de posiciones.")
	}
	for _, grp := range m.league.TablesGroups {
		for _, nt := range grp.Tables {
			lines = append(lines, m.tableLines(nt, width)...)
			lines = append(lines, "")
		}
	}
	return windowLines(lines, m.lgCursor, width, height)
}

func (m model) tableLines(nt api.NamedTable, width int) []string {
	var out []string
	if nt.Name != "" {
		out = append(out, styleSection.Render(nt.Name))
	}

	cols := nt.Table.Columns
	colW := make([]int, len(cols))
	for i, c := range cols {
		colW[i] = lipgloss.Width(c.Title) + 1
		if colW[i] < 4 {
			colW[i] = 4
		}
	}
	numW := 0
	for _, c := range colW {
		numW += c
	}
	posW := 6
	nameW := width - posW - numW - 1
	if nameW < 8 {
		nameW = 8
	}
	if nameW > 34 {
		nameW = 34
	}

	// header
	hdr := styleTblHeader.Render(padRight("  #", posW) + padRight(" Equipo", nameW))
	for i, c := range cols {
		hdr += styleTblHeader.Render(padLeft(c.Title, colW[i]))
	}
	out = append(out, hdr)

	for _, r := range nt.Table.Rows {
		zone := zoneBlock(r.DestinationColor)
		pos := padLeft(strconv.Itoa(r.Num), 3) + zone + " "
		name := styleTblName.Render(padRight(r.Entity.Object.Name, nameW))
		line := padRight(styleTblPos.Render(pos), posW) + name
		for i, c := range cols {
			val := r.Get(c.Key)
			st := styleTblCell
			if c.IsBold {
				st = styleTblPts
			}
			line += st.Render(padLeft(val, colW[i]))
		}
		out = append(out, lipgloss.NewStyle().Width(width).Render(line))
	}

	// qualification legend
	if len(nt.Table.Destinations) > 0 {
		var leg []string
		for _, d := range nt.Table.Destinations {
			leg = append(leg, zoneBlock(d.Color)+" "+styleKV.Render(d.Name))
		}
		out = append(out, strings.Join(leg, "   "))
	}
	return out
}

func (m model) renderFixtures(width, height int) string {
	gs := roundGames(m.league, m.fixRound)
	roundName := "Fixture"
	if m.fixRound < len(m.league.Games.Filters) {
		roundName = m.league.Games.Filters[m.fixRound].Name
	}
	head := styleSection.Render("◀ "+roundName+" ▶") + "\n"
	if len(gs) == 0 {
		return head + lipgloss.NewStyle().Width(width).Padding(0, 0).Foreground(colMuted).Render("Sin partidos en esta fecha.")
	}
	if m.lgCursor >= len(gs) {
		m.lgCursor = len(gs) - 1
	}

	var b strings.Builder
	b.WriteString(head)
	avail := height - 1
	start := 0
	if m.lgCursor >= avail {
		start = m.lgCursor - avail + 1
	}
	end := start + avail
	if end > len(gs) {
		end = len(gs)
	}
	for i := start; i < end; i++ {
		b.WriteString(m.renderGameRow(gs[i], width, i == m.lgCursor))
		b.WriteByte('\n')
	}
	return b.String()
}

func windowLines(lines []string, offset, width, height int) string {
	if offset < 0 {
		offset = 0
	}
	if offset > len(lines)-1 {
		offset = len(lines) - 1
	}
	if offset < 0 {
		offset = 0
	}
	end := offset + height
	if end > len(lines) {
		end = len(lines)
	}
	return lipgloss.NewStyle().Width(width).Render(strings.Join(lines[offset:end], "\n"))
}
