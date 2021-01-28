package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type playByPlay struct {
	Game struct {
		GameID  string
		Actions []struct {
			ActionNumber int
			ActionType   string
			Clock        string
			Period       int
			PeriodType   string
			TeamTricode  string
			ScoreHome    string
			ScoreAway    string
			Description  string
			Value        string
			Printed      string
		}
	}
}

type playByPlayParser struct {
	URL             string
	CurrentActionID int
	Game            Game
	PeriodMapper    map[int]string
	PlayByPlay      playByPlay
	PrintBox        string
}

func convertClock(clock string) string {
	re := regexp.MustCompile(`PT(.+)M(.+)S`)
	matches := re.FindAllStringSubmatch(clock, -1)[0]
	return fmt.Sprintf("%s:%s", matches[1], matches[2])
}

func (p *playByPlayParser) parse() {
	text := GetRequestBody(p.URL)

	err := json.Unmarshal(text, &p.PlayByPlay)

	if err != nil {
		log.Fatalf("Failed to parse gameID %s", p.Game.GameID)
	}
}

func (p *playByPlayParser) process() {
	homeName := p.Game.HomeTeam.TeamTricode
	awayName := p.Game.AwayTeam.TeamTricode

	actions := p.PlayByPlay.Game.Actions
	for i, action := range actions {
		homeScore := action.ScoreHome
		awayScore := action.ScoreAway

		description := action.Description
		if description == "" {
			description = action.Value
		}

		var team string
		if len(action.TeamTricode) > 0 {
			team = fmt.Sprintf("(%s)", action.TeamTricode)
			description = strings.Join([]string{team, description}, " ")
		}

		period := p.PeriodMapper[action.Period]
		clock := convertClock(action.Clock)
		text := fmt.Sprintf("%s %3s:%-3s %s | %s %s | %s", awayName, awayScore, homeScore, homeName, period, clock, description)
		actions[i].Printed = text
	}
}

func (p *playByPlayParser) output() {

	actions := p.PlayByPlay.Game.Actions
	for _, action := range actions {
		if action.ActionNumber > p.CurrentActionID && len(action.Description) > 0 {
			fmt.Println(action.Printed)
		}
	}

	lastAction := actions[len(actions)-1]
	if p.CurrentActionID < lastAction.ActionNumber {
		switch lastAction.ActionType {
		case "period":
			p.PrintBox = "full"
		case "game":
			p.PrintBox = "full"
		case "timeout":
			p.PrintBox = "oncourt"
		default:
			p.PrintBox = "None"
		}
	} else {
		p.PrintBox = "None"
	}

	p.CurrentActionID = actions[len(actions)-1].ActionNumber
}

func (p *playByPlayParser) update() {
	p.parse()
	p.process()
	p.output()
}

func NewPlayByPlay(game Game) playByPlayParser {
	pbp := playByPlayParser{Game: game}

	pbp.CurrentActionID = 0
	pbp.URL = fmt.Sprintf("https://cdn.nba.com/static/json/liveData/playbyplay/playbyplay_%s.json", pbp.Game.GameID)
	pbp.PeriodMapper = map[int]string{
		1: "1st",
		2: "2nd",
		3: "3rd",
		4: "4th",
		5: "ot",
		6: "2ot",
		7: "3ot",
		8: "4ot",
	}

	pbp.update()
	return pbp
}
