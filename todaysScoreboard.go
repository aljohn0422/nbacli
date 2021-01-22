package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

type team struct {
	TeamID      int
	TeamName    string
	TeamCity    string
	TeamTricode string
	Wins        int
	Losses      int
	Score       int
}

type Game struct {
	GameID         string
	GameCode       string
	GameStatus     int
	GameStatusText string
	Period         int
	GameClock      string
	GameTimeUTC    string
	GameTimeGolang time.Time
	HomeTeam       team
	AwayTeam       team
	PbOdds         struct {
		Team string
		Odds float32
	}
	Status string
}

type todaysScoreboardParser struct {
	URL              string
	TodaysScoreboard struct {
		Scoreboard struct {
			GameDate string
			Games    []Game
		}
	}
}

func (s *todaysScoreboardParser) update() {
	s.parse()
	s.process()
}

func (s *todaysScoreboardParser) parse() {
	url := s.URL
	text := GetRequestBody(url)

	err := json.Unmarshal(text, &s.TodaysScoreboard)

	if err != nil {
		panic(err)
	}
}

func (s *todaysScoreboardParser) process() {
	games := s.TodaysScoreboard.Scoreboard.Games
	for i, game := range games {
		t := toUTC(game.GameTimeUTC)
		games[i].GameTimeGolang = t
	}
	sort.Slice(games, func(i, j int) bool {
		return games[i].GameTimeGolang.Before(games[j].GameTimeGolang)
	})

	for i, game := range games {
		datetime := game.GameTimeGolang.Format("2006-01-02 15:04")

		homeTeam := game.HomeTeam
		awayTeam := game.AwayTeam

		matchup := fmt.Sprintf("%4s at %-4s", awayTeam.TeamTricode, homeTeam.TeamTricode)

		score := fmt.Sprintf("%3d:%-3d", awayTeam.Score, homeTeam.Score)
		gameStatus := game.GameStatusText

		status := fmt.Sprintf("%s %s %s %s", datetime, matchup, score, gameStatus)
		games[i].Status = status
	}
}

func (s *todaysScoreboardParser) getGameList() []string {
	list := []string{}
	for _, game := range s.TodaysScoreboard.Scoreboard.Games {
		list = append(list, game.Status)
	}
	return list
}

func toUTC(datetime string) time.Time {
	layout := "2006-01-02T15:04:05Z"
	tm, _ := time.Parse(layout, datetime)
	return tm.In(time.Now().Location())
}

func NewTodaysScoreboard() todaysScoreboardParser {
	sp := todaysScoreboardParser{URL: "https://cdn.nba.com/static/json/liveData/scoreboard/todaysScoreboard_00.json"}

	sp.update()
	return sp
}
