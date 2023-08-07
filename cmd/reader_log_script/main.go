package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/castmetal/cloudwalk-tests-quake-game-log/pkg/logger_runner"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const TIMEOUT_IN_MINUTES = 10

// GetScriptCommand - Get the script command to execute
func GetScriptCommand() *cobra.Command {
	logPath := "../../inputs/qgames.log"
	reportsPath := "../../reports_data/"

	return &cobra.Command{
		Use:   "reader_log_script",
		Short: "Read a log game file and extract all kill reports and stats",
		Run: func(cmd *cobra.Command, args []string) {
			execArg, _ := cmd.Flags().GetBool("execute")

			execute := PrintFirstTerminalInstructions(execArg)
			if !execute {
				TerminateScript()
			}

			ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_IN_MINUTES*time.Minute)
			defer cancel()

			fmt.Print("\n\n", color.HiGreenString(`Reading log file ...`))
			uiprogress.Start()

			bar := GetBarSteps()

			start := time.Now()
			runner := logger_runner.NewLoggerRunner(bar)
			runnerResponse, err := runner.Run(ctx, logPath)
			if err != nil {
				TerminateScript()
			}

			elapsed := time.Since(start)
			log.Printf("time elapsed - took %s", elapsed)

			fmt.Print("\n\n", color.HiGreenString(`Saving reports ...`), "\n\n")
			IncBar(bar)

			time.Sleep(1 * time.Second)
			start = time.Now()
			runner.SaveReports(reportsPath)
			elapsed = time.Since(start)
			log.Printf("time elapsed - took %s", elapsed)

			fmt.Print("\n", color.HiGreenString(fmt.Sprintf(`Reports saved in path: %s `, reportsPath)), "\n")
			fmt.Print("\n", color.HiGreenString(`Log was read, generating results ...`), "\n\n")
			time.Sleep(1 * time.Second)
			IncBar(bar)
			uiprogress.Stop()

			PrintReportsStats(runnerResponse)
		},
	}
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

func GetKillsByPlayerText(kills map[string]int32) string {
	killsByPlayer := ""
	for killer, numberOfKills := range kills {
		killsByPlayer += fmt.Sprintf("%s: %d \n", killer, numberOfKills)
	}

	return killsByPlayer
}

func PrintReportsStats(runnerResponse *logger_runner.RunnerResponse) {
	for i := 1; i <= runnerResponse.TotalGames; i++ {
		fmt.Print("\n", color.HiGreenString(fmt.Sprintf("Stats for game: %d - Grouped Report", i)), "\n\n")

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

		groupedReport := runnerResponse.GroupedReport[fmt.Sprintf("game-%d", i)]
		mostKillerPlayer := GetMostKillerPlayer(groupedReport.Kills)
		killsByPlayer := GetKillsByPlayerText(groupedReport.Kills)

		r := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%d", i)},
			{Align: simpletable.AlignLeft, Text: strings.Join(groupedReport.Players, ", \n")},
			{Align: simpletable.AlignLeft, Text: killsByPlayer},
			{Align: simpletable.AlignCenter, Text: mostKillerPlayer},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", groupedReport.TotalKills)},
		}

		tableGroupedReport.Body.Cells = append(tableGroupedReport.Body.Cells, r)

		tableGroupedReport.SetStyle(simpletable.StyleCompactLite)
		fmt.Println(tableGroupedReport.String())
	}

}

func IncBar(bar *uiprogress.Bar) {
	bar.Incr()
	time.Sleep(time.Millisecond * 100)
}

func GetBarSteps() *uiprogress.Bar {
	var steps = []string{"read file", "saving each kill data", "report data generated", "saving reports", "printing reports"}
	bar := uiprogress.AddBar(len(steps))

	// prepend the current step to the bar
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return "app: " + steps[b.Current()-1]
	})

	return bar
}

// PrintFirstTerminalInstructions - Print terminal instructions for dev mode
func PrintFirstTerminalInstructions(execArg bool) bool {
	if execArg {
		return true
	}

	fmt.Print("\n\n       ")
	fmt.Print(`TYPE `)
	fmt.Print(color.HiGreenString(`"y"`))
	fmt.Print(` TO EXECUTE OR `)
	fmt.Print(color.RedString(`"n"`))
	fmt.Print(` TO FINISH SCRIPT AND THEN PRESS ENTER: `)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	str := strings.TrimSpace(strings.ToLower(line))
	if str != "y" {
		return false
	}

	return true
}

// PrintPreRun - Dev mode printer
func PrintPreRun() {
	size := 1

	if term.IsTerminal(0) {
		size = 0
	}

	terminalWidth, _, err := term.GetSize(size)
	if err != nil {
		return
	}

	fmt.Print("\n\n")

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("#3C41F5")).
		Bold(true).
		Width(terminalWidth).
		Align(lipgloss.Left)

	contentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		Background(lipgloss.Color("0")).
		Width(terminalWidth).
		Padding(1)

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("#E2E2E2")).
		Bold(true).
		Width(terminalWidth).
		Align(lipgloss.Center)

	header := headerStyle.Render("Castmetal - github.com/castmetal")
	content := contentStyle.Render("Cloudwalk Reader Log Script")
	footer := footerStyle.Render("© 2023")

	fmt.Println(header)
	fmt.Println(content)
	fmt.Println(footer)

	fmt.Print("\n\n\n")
}

func TerminateScript() {
	fmt.Print("\n\n")
	fmt.Println(color.HiGreenString(`Terminando o script ...`))
	os.Exit(0)
}

func main() {
	scriptCommand := GetScriptCommand()

	scriptCommand.Flags().BoolP("execute", "e", false, "Execute script")

	scriptCommand.PreRun = func(cmd *cobra.Command, args []string) {
		execArg, _ := cmd.Flags().GetBool("execute") // If dev mode = true execute pre-run printer
		if !execArg {
			PrintPreRun()
		}
	}

	if err := scriptCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
