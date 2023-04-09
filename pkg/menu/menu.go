package menu

import (
	"os"
	"os/exec"
	"strings"
)

type Menu struct {
	Prompt  string
	Options []string
}

func Exit(errMsg string) {
	cmd := exec.Command("dunstify", errMsg, "-u", "critical")
	cmd.Run()
	os.Exit(1)
}

func ExitOnError(err error, errMsg string) {
	if err != nil {
		Exit(errMsg)
	}
}

func NotifySuccess(msg string) {
	cmd := exec.Command("dunstify", msg)
	cmd.Run()
	os.Exit(0)
}

func RunMenu(prompt string, options []string) string {
	optionsStr := strings.Join(options, "\n")

	cmd := exec.Command("dmenu", "-i", "-c", "-l", "20", "-p", prompt)
	cmd.Stdin = strings.NewReader(optionsStr)

	output, err := cmd.Output()
	ExitOnError(err, "Exit Moco menu.")

	return strings.TrimSpace(string(output))
}

func RunConfirmationMenu() {
	confirmMenuOptions := []string{"Yes", "No"}
	selected := RunMenu("Are you sure?", confirmMenuOptions)

	if selected != "Yes" {
		Exit("Operation cancelled.")
	}
}
