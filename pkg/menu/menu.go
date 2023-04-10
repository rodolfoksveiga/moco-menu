package menu

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Menu struct {
	Prompt  string
	Options []string
}

func IndexMenuOptions(menuOptions []string) []string {
	if len(menuOptions) > 10 {
		for index := range menuOptions {
			menuOptions[index] = fmt.Sprintf("%02d. %s",
				index,
				menuOptions[index],
			)
		}

		return menuOptions
	}

	for index := range menuOptions {
		menuOptions[index] = fmt.Sprintf("%s. %s",
			strconv.Itoa(index),
			menuOptions[index],
		)
	}
	return menuOptions
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

func ValidateDate(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	ExitOnError(err, "Wrong date format.")

	return date
}

func ValidateDateOutput(dateOutput string) string {
	re := regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}")
	dateStr := re.FindString(dateOutput)

	ValidateDate(dateStr)

	return dateStr
}

func ValidateDuration(selectedDuration string) float64 {
	newDuration, err := strconv.ParseFloat(selectedDuration, 64)
	ExitOnError(err, "Wrong duration format.")

	return newDuration
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
	ExitOnError(err, "Exit moco menu.")

	return strings.TrimSpace(string(output))
}

func RunConfirmationMenu() {
	confirmMenuOptions := []string{"1. Yes", "2. No"}
	selected := RunMenu("Are you sure?", confirmMenuOptions)

	if selected != "1. Yes" {
		Exit("Operation cancelled.")
	}
}
