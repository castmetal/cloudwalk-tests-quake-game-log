package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

type Game struct {
	GameNumber int
	Data       string
}

func main() {
	filePath := "./inputs/qgames.log"

	games := make(chan Game, 20)
	wg := sync.WaitGroup{}

	go processGames(games, &wg)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	gameInitRegex := regexp.MustCompile(`^.*(\d+:\d+) InitGame: .*$`)
	gameKillRegex := regexp.MustCompile(`^.*(\d+:\d+) Kill: .*$`)
	gameFinishRegex := regexp.MustCompile(`^.*(\d+:\d+) ShutdownGame:.*$`)
	gameNumber := 0
	var gameData strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		matchInitGame := gameInitRegex.FindStringSubmatch(line)
		if len(matchInitGame) > 0 {
			gameData.Reset()
			continue
		}

		matchShutdownGame := gameFinishRegex.FindStringSubmatch(line)
		if len(matchShutdownGame) > 0 {
			gameNumber++
			wg.Add(1)
			game := Game{
				GameNumber: gameNumber,
				Data:       gameData.String(),
			}
			games <- game
			gameData.Reset()
			continue
		}

		matchKillData := gameKillRegex.FindStringSubmatch(line)
		if len(matchKillData) > 0 {
			gameData.WriteString(line + "\n")
		}
	} // O(n) - n log data size

	close(games)

	wg.Wait()
}

func processGames(games <-chan Game, wg *sync.WaitGroup) {
	for game := range games {
		go func(game Game) {
			defer wg.Done()
			fmt.Printf("Game-%d:\n%s\n\n", game.GameNumber, game.Data)
		}(game)
	}
}
