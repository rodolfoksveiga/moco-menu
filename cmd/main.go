// main should be a small function that imports from /pkg and /internal
package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rodolfoksveiga/moco-menu/internal/auth"
	"github.com/rodolfoksveiga/moco-menu/pkg/api"
)

func main() {
	/* mainMenuOptions := []string{"Control activity timer", "Create activity", "Edit activity", "Delete activity"}

	selected := menu.RunMenu("Moco menu:", mainMenuOptions) */

	config := auth.Config{ConfigFilePath: "../internal/auth/moco.json"}
	authConfig, errMsg := config.Load()
	if errMsg != nil {
		cmd := exec.Command("dunstify", *errMsg, "-u", "critical")
		cmd.Run()
		os.Exit(1)
	}

	apiClient := api.Client{Domain: authConfig.Domain, Email: authConfig.Email, UserId: authConfig.UserId, AdminApiKey: authConfig.AdminApiKey, ApiKey: authConfig.ApiKey}

	activities := apiClient.FetchActivities("2023-04-08", "2023-04-08")
	fmt.Println(activities)

	apiClient.UpdateActivity(1030703750, 945824843, 12188024, "2023-04-08", 1, "test")

	/* switch selected {
	case mainMenuOptions[0]:
		timerMenuOptions := []string{"Start", "Stop"}
		selected := menu.RunMenu("Timer control:", timerMenuOptions)

		switch selected {
		case timerMenuOptions[0]:
			activities := apiClient.FetchActivities()
			fmt.Println(activities[0].Id)
		case timerMenuOptions[1]:
			fmt.Println("there")
		}
	case mainMenuOptions[1]:
		fmt.Println("Implement case")
	case mainMenuOptions[2]:
		fmt.Println("Implement case")
	case "Delete activity":
		fmt.Println("Implement case")
	default:
		cmd := exec.Command("dunstify", "Invalid menu option", "-u", "critical")
		cmd.Run()
		os.Exit(1)
	} */
}
