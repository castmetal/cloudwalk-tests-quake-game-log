package logger_runner

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"sync"
	"sync/atomic"

	"github.com/castmetal/cloudwalk-tests-quake-game-log/pkg/reports"
)

const (
	CONSUMERS       = 20
	INIT_GAME_REGEX = `^.*(\d+:\d+) InitGame: .*$`
)

type LoggerRunner struct {
	DeathsReport        reports.DeathsReport
	mu                  sync.Mutex
	TotalProcessedItems int32
	TotalItemsToProcess int32
}

type RunnerResponse struct {
	DeathsReport *reports.DeathsReport
}

type KillGameData struct {
	GameNumber int
	Data       string
}

func NewLoggerRunner() *LoggerRunner {
	return &LoggerRunner{
		DeathsReport:        make(reports.DeathsReport, 0),
		mu:                  sync.Mutex{},
		TotalProcessedItems: 0,
		TotalItemsToProcess: 0,
	}
}

func (r *LoggerRunner) Run(ctx context.Context, logPath string) (*RunnerResponse, error) {
	killData := make(chan KillGameData, CONSUMERS)
	done := make(chan bool, 1)

	// Consume each kill data with the same size of consumers by killData chan
	for i := 0; i < CONSUMERS; i++ {
		go r.processKillData(ctx, killData, done)
	}

	file, err := os.Open(logPath)
	if err != nil {
		return &RunnerResponse{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	gameInitRegex := regexp.MustCompile(INIT_GAME_REGEX)
	gameKillRegex := regexp.MustCompile(`^.*(\d+:\d+) Kill: .*$`)
	gameFinishRegex := regexp.MustCompile(`^.*(\d+:\d+) ShutdownGame:.*$`)
	gameNumber := 0
	var totalItemsToProcess int32

	for scanner.Scan() {
		line := scanner.Text()
		matchInitGame := gameInitRegex.FindStringSubmatch(line)
		if len(matchInitGame) > 0 {
			game := KillGameData{
				GameNumber: gameNumber,
				Data:       line,
			}

			killData <- game
			totalItemsToProcess++
			continue
		}

		matchShutdownGame := gameFinishRegex.FindStringSubmatch(line)
		if len(matchShutdownGame) > 0 {
			gameNumber++
			continue
		}

		matchKillData := gameKillRegex.FindStringSubmatch(line)
		if len(matchKillData) > 0 {

			game := KillGameData{
				GameNumber: gameNumber,
				Data:       line,
			}

			killData <- game

			totalItemsToProcess++
		}
	} // O(n) - n log data size

	r.mu.Lock()
	r.TotalItemsToProcess = totalItemsToProcess
	r.mu.Unlock()

	close(killData)

	for {
		select {
		case <-ctx.Done():
			return &RunnerResponse{}, fmt.Errorf("ctx_canceled")
		case <-done:
			return &RunnerResponse{
				DeathsReport: &r.DeathsReport,
			}, nil
		}
	}
}

func (r *LoggerRunner) addProcessedItems() {
	atomic.AddInt32(&r.TotalProcessedItems, 1)
}

func (r *LoggerRunner) processKillData(ctx context.Context, killData <-chan KillGameData, done chan bool) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("ctx_canceled")
		case killGameData, ok := <-killData:
			if !ok && r.TotalProcessedItems != r.TotalItemsToProcess {
				continue
			}

			if !ok && r.TotalProcessedItems == r.TotalItemsToProcess {
				done <- true
				return nil
			}

			gameStr := fmt.Sprintf("game-%d", killGameData.GameNumber)
			gameInitRegex := regexp.MustCompile(INIT_GAME_REGEX)
			matchInitGame := gameInitRegex.FindStringSubmatch(killGameData.Data)
			if len(matchInitGame) > 0 {
				mapModMeans := make(map[reports.KillMods]int32, 0) // Init kill mods
				r.DeathsReport[gameStr] = reports.Game{
					KillsByMeans: mapModMeans,
				}

				r.addProcessedItems()
				continue
			}

			fmt.Println(killGameData.Data)

			r.addProcessedItems()
		}
	}
}
