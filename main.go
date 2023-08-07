package main

import (
	"context"
	"fmt"

	"github.com/castmetal/cloudwalk-tests-quake-game-log/pkg/logger_runner"
)

func main() {
	filePath := "./inputs/qgames.log"
	ctx := context.Background()

	runner := logger_runner.NewLoggerRunner()

	runnerReponse, err := runner.Run(ctx, filePath)
	if err != nil {
		fmt.Println(err)
	}

	data := *runnerReponse.DeathsReport

	fmt.Println(data["game-1"])
}
