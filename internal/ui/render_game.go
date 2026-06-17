package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderGame(width, height int) string {
	if m.err != nil {
		return styleErr.Width(width).Render("Error: " + m.err.Error())
	}
	if m.game == nil {
		return lipgloss.NewStyle().Width(width).Padding(1, 2).Foreground(colMuted).Render(m.spinner.View() + " cargando…")
	}
	g := m.game.Game
	if len(g.Teams) < 2 {
		return lipgloss.NewStyle().Width(width).Padding(1, 2).Render("Sin datos del partido.")
	}
	home, away := g.Teams[0], g.Teams[1]

	stage := g.StageRoundName
	if g.League.Name != "" {
		stage += " · " + g.League.Name
	}
	stageLine := lipgloss.NewStyle().Foreground(colMuted).Render(stage)

	hStyle, aStyle := styleTeamWin, styleTeamWin
	if g.IsFinal() || g.HasScore() {
		if g.Winner == 0 {
			aStyle = styleTeamLose
		} else if g.Winner == 1 {
			hStyle = styleTeamLose
		}
	}

	panelW := width - 4
	innerW := panelW - 2
	colW := (innerW - 9) / 2
	if colW < 12 {
		colW = 12
	}
	hCell := lipgloss.NewStyle().Width(colW).Align(lipgloss.Right).Render(
		hStyle.Render(home.Name) + " " + teamColorBlock(home.Colors.Color))
	aCell := lipgloss.NewStyle().Width(colW).Align(lipgloss.Left).Render(
		teamColorBlock(away.Colors.Color) + " " + aStyle.Render(away.Name))

	scoreStyle := styleScore
	if g.IsLive() {
		scoreStyle = styleScoreLive
	} else if !g.HasScore() {
		scoreStyle = styleClock
	}
	score := scoreStyle.Bold(true).Render(" " + g.ScoreText() + " ")
	scoreRow := lipgloss.JoinHorizontal(lipgloss.Center, hCell, score, aCell)

	var statusS string
	switch {
	case g.IsLive():
		statusS = styleClockLive.Render("● " + g.Clock())
	case g.IsFinal():
		statusS = styleClockFinal.Render(g.Clock())
	default:
		statusS = styleClock.Render(g.Clock())
	}
	statusRow := lipgloss.NewStyle().Width(innerW).Align(lipgloss.Center).Render(statusS)

	card := stylePanel.Width(panelW).Render(
		lipgloss.JoinVertical(lipgloss.Center, scoreRow, statusRow))

	innerLeft := width - 2

	var sections []string
	sections = append(sections, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(stageLine))
	sections = append(sections, card)

	if goals := m.renderGoals(g); goals != "" {
		sections = append(sections, "", styleSection.Render("⚽ Goles"), goals)
	}

	if cards := m.renderCards(g); cards != "" {
		sections = append(sections, "", styleSection.Render("🟨 Tarjetas"), cards)
	}

	if len(g.TVNetworks) > 0 {
		var tv []string
		for _, n := range g.TVNetworks {
			tv = append(tv, n.Name)
		}
		sections = append(sections, styleKV.Render("📺 TV: ")+styleKVval.Render(strings.Join(tv, ", ")))
	}

	if stats := m.renderStats(g, innerLeft); stats != "" {
		sections = append(sections, "", styleSection.Render("📊 Estadísticas"), stats)
	}

	if len(g.GameInfo) > 0 {
		sections = append(sections, "")
		sections = append(sections, styleSection.Render("Información"))
		for _, it := range g.GameInfo {
			sections = append(sections, styleKV.Render(" "+it.Name+": ")+styleKVval.Render(it.Value))
		}
	}

	body := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return lipgloss.NewStyle().Width(width).MaxHeight(height).Padding(1, 1).Render(body)
}
