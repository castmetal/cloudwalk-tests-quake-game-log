package console_helpers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/gosuri/uiprogress"
	"golang.org/x/term"
)

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
