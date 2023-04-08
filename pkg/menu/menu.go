package menu

import (
	"os"
	"os/exec"
	"strings"
)

func ExitOnError(err error, msg string) {
	if err != nil {
		cmd := exec.Command("dunstify", msg, "-u", "critical")
		cmd.Run()
		os.Exit(1)
	}
}

func RunMenu(prompt string, options []string) string {
	optionsStr := strings.Join(options, "\n")

	cmd := exec.Command("dmenu", "-i", "-c", "-l", "20", "-p", prompt)
	cmd.Stdin = strings.NewReader(optionsStr)

	output, err := cmd.Output()
	ExitOnError(err, "Exit Menu")

	return strings.TrimSpace(string(output))
}
