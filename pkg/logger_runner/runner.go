package logger_runner

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/castmetal/cloudwalk-tests-quake-game-log/pkg/reports"
)

const (
	CONSUMERS         = 20
	INIT_GAME_REGEX   = `^.*(\d+:\d+) InitGame:.*$`
	KILL_REGEX        = `^.*(\d+:\d+) Kill: .*$`
	KILL_REMOVE_REGEX = `^\d+:\d+ Kill: \d+ \d+ \d+: `
	WORLD_PLAYER_ID   = `<world>`
)

type LoggerRunner struct {
	DeathsReport        reports.DeathsReport
	GroupedReport       reports.GroupedReport
	mu                  sync.Mutex
	TotalProcessedItems int32
	TotalItemsToProcess int32
	PlayerData          map[string]map[string]bool
}

type RunnerResponse struct {
	DeathsReport  *reports.DeathsReport
	GroupedReport *reports.GroupedReport
}

type KillGameData struct {
	GameNumber int
	Data       string
}

type PlayerKilledData struct {
	KillerPlayer string
	Mod          reports.KillMods
	DeadPlayer   string
}

// NewLoggerRunner - Get a LoggerRunner to Execute
func NewLoggerRunner() *LoggerRunner {
	return &LoggerRunner{
		DeathsReport:        make(reports.DeathsReport, 0),
		GroupedReport:       make(reports.GroupedReport, 0),
		mu:                  sync.Mutex{},
		TotalProcessedItems: 0,
		TotalItemsToProcess: 0,
		PlayerData:          make(map[string]map[string]bool, 0),
	}
}

// Run - Run a logger and get a response containing AllReports - T O(n) - n log data size - S O (n*m)
//   - Input - ctx:  context.Context
//   - Input - ctx:  logPath string
//   - Response - RunnerResponse:  response struct with all generated reports
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
	gameKillRegex := regexp.MustCompile(KILL_REGEX)
	gameInitRegex := regexp.MustCompile(INIT_GAME_REGEX)
	gameNumber := 0
	var totalItemsToProcess int32

	for scanner.Scan() {
		line := scanner.Text()

		matchShutdownGame := gameInitRegex.FindStringSubmatch(line)
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
	}

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
				DeathsReport:  &r.DeathsReport,
				GroupedReport: &r.GroupedReport,
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
			// Consume till TotalProcessedItems = TotalItemsToProcess
			if !ok && r.TotalProcessedItems != r.TotalItemsToProcess {
				continue
			}

			if !ok && r.TotalProcessedItems == r.TotalItemsToProcess {
				done <- true
				return nil
			}

			gameStr := fmt.Sprintf("game-%d", killGameData.GameNumber)
			r.initReportMaps(gameStr)

			playerKilledData := r.getPlayerKilledData(killGameData.Data)

			r.mu.Lock()
			r.DeathsReport[gameStr].KillsByMeans[playerKilledData.Mod]++

			if r.PlayerData[gameStr] == nil {
				mapPlayerFoundedByGame := make(map[string]bool, 0)
				r.PlayerData[gameStr] = mapPlayerFoundedByGame
			}

			r.assignPlayerToList(gameStr, playerKilledData.KillerPlayer)
			r.assignPlayerToList(gameStr, playerKilledData.DeadPlayer)

			groupedEntry := r.GroupedReport[gameStr]
			groupedEntry.TotalKills++

			if playerKilledData.KillerPlayer != WORLD_PLAYER_ID {
				groupedEntry.Kills[playerKilledData.KillerPlayer]++
			} else {
				groupedEntry.Kills[playerKilledData.DeadPlayer]--
			}

			r.GroupedReport[gameStr] = groupedEntry

			r.mu.Unlock()
			r.addProcessedItems()
		}
	}
}

func (r *LoggerRunner) assignPlayerToList(gameStr string, playerId string) {
	if playerId == WORLD_PLAYER_ID {
		return
	}

	groupedEntry := r.GroupedReport[gameStr]
	if r.PlayerData[gameStr][playerId] == false {
		r.PlayerData[gameStr][playerId] = true
		groupedEntry.Players = append(groupedEntry.Players, playerId)
	}

	if groupedEntry.Kills[playerId] == 0 {
		groupedEntry.Kills[playerId] = 0
	}

	r.GroupedReport[gameStr] = groupedEntry
}

func (r *LoggerRunner) initReportMaps(gameStr string) {
	r.mu.Lock()
	if r.DeathsReport[gameStr].KillsByMeans == nil {
		mapModMeans := make(map[reports.KillMods]int32, 0) // Init kill mods
		r.DeathsReport[gameStr] = reports.Game{
			KillsByMeans: mapModMeans,
		}
		mapKills := make(map[string]int32, 0) // Init kills
		r.GroupedReport[gameStr] = reports.GroupedInformationReport{
			Kills:   mapKills,
			Players: make([]string, 0),
		}
	}
	r.mu.Unlock()
}

func (r *LoggerRunner) getPlayerKilledData(killLogLine string) *PlayerKilledData {
	str := strings.TrimSpace(killLogLine)
	regex := regexp.MustCompile(KILL_REMOVE_REGEX)
	replacedString := regex.ReplaceAllString(str, "")

	divideStr := strings.Split(replacedString, " ")

	killerPlayer := ""
	deadPlayer := ""
	foundKiller := false
	for i := 0; i < len(divideStr)-2; i++ {
		if divideStr[i] == "killed" {
			foundKiller = true
			continue
		}

		if foundKiller == false {
			killerPlayer += " " + divideStr[i]
			continue
		}

		deadPlayer += " " + divideStr[i]
	}

	var killMod reports.KillMods
	killMod = killMod.GetModByString(divideStr[len(divideStr)-1])

	return &PlayerKilledData{
		Mod:          killMod,
		KillerPlayer: strings.TrimSpace(killerPlayer),
		DeadPlayer:   strings.TrimSpace(deadPlayer),
	}
}
