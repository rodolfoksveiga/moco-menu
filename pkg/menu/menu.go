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

func Exit(msg string) {
	cmd := exec.Command("dunstify", msg, "-u", "critical")
	cmd.Run()
	os.Exit(1)
}

func ExitOnError(err error, msg string) {
	if err != nil {
		Exit(msg)
	}
}

func (menu Menu) RunMenu() string {
	optionsStr := strings.Join(menu.Options, "\n")

	cmd := exec.Command("dmenu", "-i", "-c", "-l", "20", "-p", menu.Prompt)
	cmd.Stdin = strings.NewReader(optionsStr)

	output, err := cmd.Output()
	ExitOnError(err, "Exit Menu")

	return strings.TrimSpace(string(output))
}
