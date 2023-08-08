package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/castmetal/cloudwalk-tests-quake-game-log/cmd/reader_log_script/console_helpers"
	"github.com/castmetal/cloudwalk-tests-quake-game-log/pkg/logger_runner"
	"github.com/fatih/color"
	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
)

const TIMEOUT_IN_MINUTES = 10

func main() {
	scriptCommand := GetScriptCommand()

	scriptCommand.Flags().BoolP("execute", "e", false, "Execute script")

	scriptCommand.PreRun = func(cmd *cobra.Command, args []string) {
		execArg, _ := cmd.Flags().GetBool("execute") // If dev mode = true execute pre-run printer
		if !execArg {
			console_helpers.PrintPreRun()
		}
	}

	if err := scriptCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func GetAbsPath() string {
	path, err := os.Getwd()
	if err != nil {
		return "/root"
	}

	if strings.Contains(path, "/cmd") {
		splitPath := strings.Split(path, "/cmd")
		path = splitPath[0]
	}

	return path
}

// GetScriptCommand - Get the script command to execute
func GetScriptCommand() *cobra.Command {
	path := GetAbsPath()
	logPath := path + "/inputs/qgames.log"
	reportsPath := path + "/reports_data/"

	return &cobra.Command{
		Use:   "reader_log_script",
		Short: "Read a log game file and extract all kill reports and stats",
		Run: func(cmd *cobra.Command, args []string) {
			execArg, _ := cmd.Flags().GetBool("execute")

			execute := console_helpers.PrintFirstTerminalInstructions(execArg)
			if !execute {
				console_helpers.TerminateScript()
			}

			ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_IN_MINUTES*time.Minute)
			defer cancel()

			fmt.Print("\n\n", color.HiGreenString(`Reading log file ...`))
			uiprogress.Start()

			bar := console_helpers.GetBarSteps()

			start := time.Now()
			runner := logger_runner.NewLoggerRunner(bar)
			runnerResponse, err := runner.Run(ctx, logPath)
			if err != nil {
				console_helpers.TerminateScript()
			}

			elapsed := time.Since(start)
			log.Printf("time elapsed - took %s", elapsed)

			fmt.Print("\n\n", color.HiGreenString(`Saving reports ...`), "\n\n")
			console_helpers.IncBar(bar)

			time.Sleep(1 * time.Second)
			start = time.Now()
			runner.SaveReports(reportsPath)
			elapsed = time.Since(start)
			log.Printf("time elapsed - took %s", elapsed)

			fmt.Print("\n", color.HiGreenString(fmt.Sprintf(`Reports saved in path: %s `, reportsPath)), "\n")
			fmt.Print("\n", color.HiGreenString(`Log was read, generating results ...`), "\n\n")
			time.Sleep(1 * time.Second)
			console_helpers.IncBar(bar)
			uiprogress.Stop()

			console_helpers.PrintReportsStats(runnerResponse)
		},
	}
}
