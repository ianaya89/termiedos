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

type Team struct {
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	URLName   string `json:"url_name"`
	ID        string `json:"id"`
	CountryID string `json:"country_id"`
	Colors    Colors `json:"colors"`
	RedCards  int    `json:"red_cards"`
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

type GameCenter struct {
	TTL  int        `json:"TTL"`
	Game GameDetail `json:"game"`
}

type GameDetail struct {
	Game
	League   League         `json:"league"`
	GameInfo []GameInfoItem `json:"game_info"`
}
