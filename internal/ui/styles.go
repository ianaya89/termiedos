package ui

import "github.com/charmbracelet/lipgloss"

// Palette tuned to mimic the promiedos.com.ar dark theme.
var (
	colBg     = lipgloss.Color("#0D0E12")
	colPanel  = lipgloss.Color("#161821")
	colHeader = lipgloss.Color("#1F2330")
	colBorder = lipgloss.Color("#2A2F3C")
	colText   = lipgloss.Color("#E8EAED")
	colMuted  = lipgloss.Color("#8B909C")
	colDim    = lipgloss.Color("#5A5F6B")
	colLive   = lipgloss.Color("#22D366")
	colAccent = lipgloss.Color("#19E54B")
	colSelBg  = lipgloss.Color("#2D3344")
	colScore  = lipgloss.Color("#FFFFFF")
	colRed    = lipgloss.Color("#E5484D")
)

var (
	styleTitle = lipgloss.NewStyle().
			Bold(true).Foreground(colBg).Background(colAccent).Padding(0, 1)

	styleTitleDate = lipgloss.NewStyle().
			Foreground(colText).Background(colHeader).Padding(0, 1)

	styleHelp = lipgloss.NewStyle().Foreground(colDim).Padding(0, 1)

	styleErr = lipgloss.NewStyle().Foreground(colRed).Bold(true).Padding(0, 1)

	styleLeagueHeader = lipgloss.NewStyle().
				Bold(true).Foreground(colText).Background(colHeader).Padding(0, 1)

	styleLeagueCountry = lipgloss.NewStyle().Foreground(colMuted)

	styleClock      = lipgloss.NewStyle().Foreground(colMuted)
	styleClockLive  = lipgloss.NewStyle().Foreground(colLive).Bold(true)
	styleClockFinal = lipgloss.NewStyle().Foreground(colDim)

	styleTeam      = lipgloss.NewStyle().Foreground(colText)
	styleTeamWin   = lipgloss.NewStyle().Foreground(colText).Bold(true)
	styleTeamLose  = lipgloss.NewStyle().Foreground(colMuted)
	styleScore     = lipgloss.NewStyle().Foreground(colScore).Bold(true)
	styleScoreLive = lipgloss.NewStyle().Foreground(colLive).Bold(true)

	styleRowSel = lipgloss.NewStyle().Background(colSelBg)

	styleSidebarTitle = lipgloss.NewStyle().Bold(true).Foreground(colMuted).Padding(0, 1)
	styleSidebarItem  = lipgloss.NewStyle().Foreground(colText).Padding(0, 1)
	styleSidebarSel   = lipgloss.NewStyle().Foreground(colAccent).Bold(true).Background(colSelBg).Padding(0, 1)

	styleTabActive   = lipgloss.NewStyle().Bold(true).Foreground(colBg).Background(colAccent).Padding(0, 2)
	styleTabInactive = lipgloss.NewStyle().Foreground(colMuted).Background(colPanel).Padding(0, 2)

	styleTblHeader = lipgloss.NewStyle().Bold(true).Foreground(colMuted)
	styleTblPos    = lipgloss.NewStyle().Foreground(colMuted)
	styleTblName   = lipgloss.NewStyle().Foreground(colText)
	styleTblPts    = lipgloss.NewStyle().Bold(true).Foreground(colText)
	styleTblCell   = lipgloss.NewStyle().Foreground(colMuted)

	stylePanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).BorderForeground(colBorder)

	styleSection = lipgloss.NewStyle().Bold(true).Foreground(colAccent)
	styleKV      = lipgloss.NewStyle().Foreground(colMuted)
	styleKVval   = lipgloss.NewStyle().Foreground(colText)

	styleEventTime = lipgloss.NewStyle().Foreground(colDim)
	styleGoal      = lipgloss.NewStyle().Foreground(colText)
	styleYellow    = lipgloss.NewStyle().Foreground(lipgloss.Color("#E5C100"))
	styleRed       = lipgloss.NewStyle().Foreground(colRed)
	styleBarHome   = lipgloss.NewStyle().Foreground(colAccent)
	styleBarAway   = lipgloss.NewStyle().Foreground(colDim)
)

var yellowMark = lipgloss.NewStyle().Foreground(lipgloss.Color("#E5C100")).Render("▮")
var redMark = lipgloss.NewStyle().Foreground(colRed).Render("▮")

func cardMark(eventType int) string {
	switch eventType {
	case 5: // straight red
		return redMark
	case 6: // second yellow → red
		return yellowMark + redMark
	default: // yellow
		return yellowMark
	}
}

func teamColorBlock(hex string) string {
	if hex == "" {
		hex = "#444444"
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(hex)).Render("▌")
}

func zoneBlock(hex string) string {
	if hex == "" {
		return " "
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(hex)).Render("▌")
}
