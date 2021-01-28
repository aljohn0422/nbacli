package main

import (
	"fmt"
	"time"

	"github.com/manifoldco/promptui"
)

func main() {
	scoreboard := NewTodaysScoreboard()

	prompt := promptui.Select{
		Label: fmt.Sprintf("Select Game"),
		Items: scoreboard.getGameList(),
		Size:  15,
	}

	index, _, err := prompt.Run()

	if err != nil {
		panic(err)
	}

	games := scoreboard.TodaysScoreboard.Scoreboard.Games
	pbp := NewPlayByPlay(games[index])
	box := NewBoxscore(games[index].GameID)
	for {
		if pbp.PrintBox == "full" {
			box.update(true)
		} else if pbp.PrintBox == "oncourt" {
			box.update(false)
		}
		time.Sleep(time.Second * 10)
		pbp.update()
	}
}
