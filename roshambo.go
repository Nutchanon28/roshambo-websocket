package main

import "fmt"

// input: player1Move, player2Move
// output: "playerX won"
func Roshambo(firstPlayerMove string, secondPlayerMove string) string {
	moveWinAgainstMap := make(map[string]string)

	moveWinAgainstMap["scissor"] = "paper"
	moveWinAgainstMap["paper"] = "rock"
	moveWinAgainstMap["rock"] = "scissor"

	firstPlayerWinAgainst, firstInputValid := moveWinAgainstMap[firstPlayerMove]
	secondPlayerWinAgainst, secondInputValid := moveWinAgainstMap[secondPlayerMove]

	if !firstInputValid || !secondInputValid {
		return "Invalid moves detected!"
	}

	var winner string
	if firstPlayerWinAgainst == secondPlayerMove {
		winner = "first player"
	} else if secondPlayerWinAgainst == firstPlayerMove {
		winner = "second player"
	} else {
		return "Draw..."
	}

	return fmt.Sprintf("The %s won!", winner)
}
