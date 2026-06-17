package api

import "strconv"

// Status enums observed from the promiedos API.
const (
	StatusScheduled = 1
	StatusLive      = 2
	StatusFinal     = 3
)

type Colors struct {
	Color     string `json:"color"`
	TextColor string `json:"text_color"`
}

type Goal struct {
	PlayerName    string  `json:"player_name"`
	PlayerSName   string  `json:"player_sname"`
	Time          float64 `json:"time"`
	TimeToDisplay string  `json:"time_to_display"`
}

type Team struct {
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	URLName   string `json:"url_name"`
	ID        string `json:"id"`
	CountryID string `json:"country_id"`
	Colors    Colors `json:"colors"`
	RedCards  int    `json:"red_cards"`
	Goals     []Goal `json:"goals"`
}

type Status struct {
	Enum       int    `json:"enum"`
	Name       string `json:"name"`
	ShortName  string `json:"short_name"`
	SymbolName string `json:"symbol_name"`
}

type TVNetwork struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Game struct {
	ID                      string      `json:"id"`
	StageRoundName          string      `json:"stage_round_name"`
	Winner                  int         `json:"winner"`
	Teams                   []Team      `json:"teams"`
	Scores                  []float64   `json:"scores"`
	URLName                 string      `json:"url_name"`
	Status                  Status      `json:"status"`
	StartTime               string      `json:"start_time"`
	GameTime                int         `json:"game_time"`
	GameTimeToDisplay       string      `json:"game_time_to_display"`
	GameTimeStatusToDisplay string      `json:"game_time_status_to_display"`
	TVNetworks              []TVNetwork `json:"tv_networks"`
}

func (g Game) IsLive() bool   { return g.Status.Enum == StatusLive }
func (g Game) IsFinal() bool  { return g.Status.Enum == StatusFinal }
func (g Game) HasScore() bool { return len(g.Scores) >= 2 }

// ScoreText renders the home-away score, e.g. "2-0", or "-" when unplayed.
func (g Game) ScoreText() string {
	if !g.HasScore() {
		return "-"
	}
	return itoa(g.Scores[0]) + "-" + itoa(g.Scores[1])
}

// Clock returns the time/minute label shown on the row.
func (g Game) Clock() string {
	if g.IsLive() && g.GameTimeToDisplay != "" {
		return g.GameTimeToDisplay
	}
	if g.GameTimeStatusToDisplay != "" {
		return g.GameTimeStatusToDisplay
	}
	if len(g.StartTime) >= 16 {
		return g.StartTime[11:16]
	}
	return g.Status.ShortName
}

func itoa(f float64) string { return strconv.Itoa(int(f)) }

type League struct {
	Name            string `json:"name"`
	ID              string `json:"id"`
	URLName         string `json:"url_name"`
	CountryID       string `json:"country_id"`
	ShowCountryFlag bool   `json:"show_country_flags"`
	CountryName     string `json:"country_name"`
	IsInternational bool   `json:"is_international"`
	Games           []Game `json:"games"`
}

type GamesResponse struct {
	Leagues   []League `json:"leagues"`
	TTL       int      `json:"TTL"`
	CacheTime int      `json:"cache_time"`
}

// ---- league standings + fixtures ----

type Column struct {
	Key    string `json:"key"`
	Title  string `json:"title"`
	Type   int    `json:"type"`
	IsBold bool   `json:"is_bold"`
}

type Destination struct {
	Num   int    `json:"num"`
	Color string `json:"color"`
	Name  string `json:"name"`
}

type CellValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Entity struct {
	Type   int  `json:"type"`
	Object Team `json:"object"`
}

type Row struct {
	Num              int         `json:"num"`
	Values           []CellValue `json:"values"`
	Entity           Entity      `json:"entity"`
	DestinationColor string      `json:"destination_color"`
}

func (r Row) Get(key string) string {
	for _, v := range r.Values {
		if v.Key == key {
			return v.Value
		}
	}
	return ""
}

type Table struct {
	IsLive       bool          `json:"is_live"`
	Destinations []Destination `json:"destinations"`
	Columns      []Column      `json:"columns"`
	Rows         []Row         `json:"rows"`
}

type NamedTable struct {
	Name  string `json:"name"`
	Table Table  `json:"table"`
}

type TableGroup struct {
	Name   string       `json:"name"`
	Tables []NamedTable `json:"tables"`
}

type FixtureFilter struct {
	Name     string `json:"name"`
	Key      string `json:"key"`
	Selected bool   `json:"selected"`
	Games    []Game `json:"games"`
}

type Fixtures struct {
	Filters []FixtureFilter `json:"filters"`
}

type LeagueData struct {
	TTL          int          `json:"TTL"`
	League       League       `json:"league"`
	TablesGroups []TableGroup `json:"tables_groups"`
	Games        Fixtures     `json:"games"`
}

// ---- gamecenter ----

type GameInfoItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Event types observed in the gamecenter timeline.
const (
	EventGoal         = 1
	EventOwnGoal      = 2
	EventGoalAlt      = 3
	EventYellowCard   = 4
	EventRedCard      = 5
	EventSecondYellow = 6
	EventPenaltyGoal  = 7
	EventSubstitution = 15
)

type EventItem struct {
	Type      int      `json:"type"`
	Time      string   `json:"time"`
	Team      int      `json:"team"` // 1 = home, 2 = away
	Texts     []string `json:"texts"`
	JerseyNum int      `json:"player_jersey_num"`
}

func (e EventItem) IsCard() bool {
	switch e.Type {
	case EventYellowCard, EventRedCard, EventSecondYellow:
		return true
	}
	return false
}

type EventRow struct {
	Time   string      `json:"time"`
	Events []EventItem `json:"events"`
}

type EventStage struct {
	Name             string     `json:"name"`
	ShowStageTitle   bool       `json:"show_stage_title"`
	IsPenaltiesStage bool       `json:"is_penalties_stage"`
	Scores           []float64  `json:"scores"`
	Rows             []EventRow `json:"rows"`
}

type StatItem struct {
	Name        string    `json:"name"`
	Values      []string  `json:"values"`
	Percentages []float64 `json:"percentages"`
}

type GameCenter struct {
	TTL  int        `json:"TTL"`
	Game GameDetail `json:"game"`
}

type GameDetail struct {
	Game
	League     League         `json:"league"`
	GameInfo   []GameInfoItem `json:"game_info"`
	Events     []EventStage   `json:"events"`
	Statistics []StatItem     `json:"statistics"`
}

// Cards returns booking events (yellow/red) in chronological order.
func (g GameDetail) Cards() []EventItem {
	var out []EventItem
	for _, st := range g.Events {
		if st.IsPenaltiesStage {
			continue
		}
		for _, r := range st.Rows {
			for _, e := range r.Events {
				if e.IsCard() {
					out = append(out, e)
				}
			}
		}
	}
	return out
}
