package console_helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/castmetal/cloudwalk-tests-quake-game-log/pkg/logger_runner"
	"github.com/castmetal/cloudwalk-tests-quake-game-log/pkg/reports"
	"github.com/fatih/color"
)

type Stats interface {
	map[string]int32 | reports.KillModMeans
}

type StatData[T Stats] struct {
	Data T
}

func GetMostKillerPlayer(kills map[string]int32) string {
	var higherKillStat int32 = -99999
	var mostKillerPlayer string
	for killer, numberOfKills := range kills {
		if numberOfKills > higherKillStat {
			mostKillerPlayer = killer
			higherKillStat = numberOfKills
		}
	}

	return mostKillerPlayer
}

func GetMostKillMean(killByMeans map[reports.KillMods]int32) reports.KillMods {
	var higherKillStat int32 = -99999
	var mostKillerMod reports.KillMods
	for mod, numberOfKills := range killByMeans {
		if numberOfKills > higherKillStat {
			mostKillerMod = mod
			higherKillStat = numberOfKills
		}
	}

	return mostKillerMod
}

func GetTotalKillsByMeans(killByMeans map[reports.KillMods]int32) int32 {
	var total int32 = 0

	for _, numberOfKills := range killByMeans {
		total += numberOfKills
	}

	return total
}

func DivideStats[T Stats](data StatData[T]) string {
	statsStr := ""

	switch v := any(data.Data).(type) {
	case map[string]int32:

		for statStr, statInt := range v {
			statsStr += fmt.Sprintf("%s: %d \n", statStr, statInt)
		}

	case reports.KillModMeans:
		for statStr, statInt := range v {
			statsStr += fmt.Sprintf("%s: %d \n", statStr.GetStrModByType(), statInt)
		}
	}

	return statsStr
}

func PrintReportsStats(runnerResponse *logger_runner.RunnerResponse) {
	for i := 1; i <= runnerResponse.TotalGames; i++ {
		PrintGroupedStats(i, runnerResponse)
		time.Sleep(30 * time.Millisecond)
	}

	for i := 1; i <= runnerResponse.TotalGames; i++ {
		PrintDeathsStats(i, runnerResponse)
		time.Sleep(30 * time.Millisecond)
	}
}

func PrintGroupedStats(gameNumber int, runnerResponse *logger_runner.RunnerResponse) {
	fmt.Print("\n", color.HiGreenString(fmt.Sprintf("Stats for game: %d - Grouped Report", gameNumber)), "\n\n")

	tableGroupedReport := simpletable.New()

	tableGroupedReport.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Game Number"},
			{Align: simpletable.AlignCenter, Text: "Players"},
			{Align: simpletable.AlignCenter, Text: "Kills"},
			{Align: simpletable.AlignCenter, Text: "Most Killer Player"},
			{Align: simpletable.AlignCenter, Text: "Total Kills"},
		},
	}

	groupedReport := runnerResponse.GroupedReport[fmt.Sprintf("game-%d", gameNumber)]
	mostKillerPlayer := GetMostKillerPlayer(groupedReport.Kills)
	statData := StatData[map[string]int32]{
		Data: groupedReport.Kills,
	}

	killsByPlayer := DivideStats(statData)

	r := []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%d", gameNumber)},
		{Align: simpletable.AlignLeft, Text: strings.Join(groupedReport.Players, ", \n")},
		{Align: simpletable.AlignLeft, Text: killsByPlayer},
		{Align: simpletable.AlignCenter, Text: mostKillerPlayer},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", groupedReport.TotalKills)},
	}

	tableGroupedReport.Body.Cells = append(tableGroupedReport.Body.Cells, r)

	tableGroupedReport.SetStyle(simpletable.StyleRounded)
	fmt.Println(tableGroupedReport.String())
}

func PrintDeathsStats(gameNumber int, runnerResponse *logger_runner.RunnerResponse) {
	fmt.Print("\n", color.HiGreenString(fmt.Sprintf("Stats for game: %d - Kills Report by Mod", gameNumber)), "\n\n")

	tableGroupedReport := simpletable.New()

	tableGroupedReport.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Game Number"},
			{Align: simpletable.AlignCenter, Text: "Kills by Means"},
			{Align: simpletable.AlignCenter, Text: "Most Kill Mean"},
			{Align: simpletable.AlignCenter, Text: "Total Kills"},
		},
	}

	deathsReport := runnerResponse.DeathsReport[fmt.Sprintf("game-%d", gameNumber)]

	statData := StatData[reports.KillModMeans]{
		Data: deathsReport.KillsByMeans,
	}

	killsByMeans := DivideStats(statData)
	mostKillMod := GetMostKillMean(deathsReport.KillsByMeans)
	totalKills := GetTotalKillsByMeans(deathsReport.KillsByMeans)

	r := []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%d", gameNumber)},
		{Align: simpletable.AlignLeft, Text: killsByMeans},
		{Align: simpletable.AlignCenter, Text: mostKillMod.GetStrModByType()},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", totalKills)},
	}

	tableGroupedReport.Body.Cells = append(tableGroupedReport.Body.Cells, r)

	tableGroupedReport.SetStyle(simpletable.StyleRounded)
	fmt.Println(tableGroupedReport.String())
}
