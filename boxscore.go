package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
)

type boxscore struct {
	Game struct {
		GameID   string
		Duration int
		HomeTeam teamBox
		AwayTeam teamBox
	}
}

type teamBox struct {
	TeamCity    string
	TeamName    string
	TeamTricode string
	Periods     []struct {
		Period int
		Score  int
	}
	Players    []player
	Statistics statistics
}

type player struct {
	Name       string
	NameI      string
	FirstName  string
	FamilyName string
	JerseyNum  string
	Position   string
	Starter    string
	Oncourt    string
	Played     string
	Statistics statistics
}

type statistics struct {
	Assists                 int
	Blocks                  int
	BlocksReceived          int
	FieldGoalsAttempted     int
	FieldGoalsMade          int
	FieldGoalsPercentage    float32
	FoulsOffensive          int
	FoulsDrawn              int
	FoulsPersonal           int
	FoulsTeam               int
	FoulsTechnical          int
	FoulsTeamTechnical      int
	FreeThrowsAttempted     int
	FreeThrowsMade          int
	FreeThrowsPercentage    float32
	Minutes                 string
	PlusMinusPoints         float32
	Points                  int
	PointsFastBreak         int
	PointsInThePaint        int
	PointsSecondChance      int
	ReboundsDefensive       int
	ReboundsOffensive       int
	ReboundsTotal           int
	Steals                  int
	ThreePointersAttempted  int
	ThreePointersMade       int
	ThreePointersPercentage float32
	Turnovers               int
	TwoPointersAttempted    int
	TwoPointersMade         int
	TwoPointersPercentage   float32
}

type boxscoreHandler struct {
	URL      string
	Boxscore boxscore
}

func (b *boxscoreHandler) parse() {
	text := GetRequestBody(b.URL)
	err := json.Unmarshal(text, &b.Boxscore)

	if err != nil {
		log.Fatalf("Failed to parse boxscore, %s", err)
	}
}

//
// Arguments
// full (bool)
// if true, print full boxscore; otherwise, print players on court.
func (b *boxscoreHandler) output(full bool) {
	fmt.Println()
	for _, t := range []teamBox{b.Boxscore.Game.AwayTeam, b.Boxscore.Game.HomeTeam} {
		fmt.Printf("%s %s\n", t.TeamCity, t.TeamName)
		fmt.Printf("%-20s%2s %5s %5s %5s %5s %3s %3s %3s %3s %3s %3s %3s %3s %3s %3s\n", "NAME", "", "TIME", "FG", "3PT", "FT", "RO", "RD", "TR", "AST", "STL", "BLK", "TO", "PF", "PTS", "+/-")
		for _, p := range t.Players {
			if !full {
				if p.Oncourt == "0" {
					continue
				}
			}
			minute := playTime(p.Statistics.Minutes)
			fmt.Printf("%-20s%2s %5s %s\n", p.NameI, p.Position, minute, stats(p.Statistics))
		}
		fmt.Printf("%29s%s\n", "", stats(t.Statistics))
	}
	fmt.Println()
}

func (b *boxscoreHandler) update(full bool) {
	b.parse()
	b.output(full)
}

func stats(s statistics) string {
	fg := fmt.Sprintf("%d-%d", s.FieldGoalsMade, s.FieldGoalsAttempted)
	thrPts := fmt.Sprintf("%d-%d", s.ThreePointersMade, s.ThreePointersAttempted)
	fthrows := fmt.Sprintf("%d-%d", s.FreeThrowsMade, s.FreeThrowsAttempted)
	val := fmt.Sprintf("%5s %5s %5s %3d %3d %3d %3d %3d %3d %3d %3d %3d %3d", fg, thrPts, fthrows, s.ReboundsOffensive, s.ReboundsDefensive, s.ReboundsTotal, s.Assists, s.Steals, s.Blocks, s.Turnovers, s.FoulsPersonal, s.Points, int(s.PlusMinusPoints))
	return val
}

func playTime(minute string) string {
	re := regexp.MustCompile(`PT(\d+)M(\d+)`)
	matches := re.FindAllStringSubmatch(minute, -1)[0]
	return fmt.Sprintf("%s:%s", matches[1], matches[2])
}

func NewBoxscore(gameID string) boxscoreHandler {
	url := fmt.Sprintf("https://cdn.nba.com/static/json/liveData/boxscore/boxscore_%s.json", gameID)
	b := boxscoreHandler{URL: url}
	return b
}
