// main should be a small function that imports from /pkg and /internal
package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"

	"github.com/rodolfoksveiga/moco-menu/internal/config"
	"github.com/rodolfoksveiga/moco-menu/pkg/api"
)

func main() {
	/* mainMenuOptions := []string{"Control activity timer", "Create activity", "Edit activity", "Delete activity"}

	selected := menu.RunMenu("Moco menu:", mainMenuOptions) */

	osUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	homeDir := osUser.HomeDir

	config := config.Config{ConfigFilePath: fmt.Sprintf("%s/.config/moco-menu/config.json", homeDir)}

	authConfig, errMsg := config.Load()
	if errMsg != nil {
		cmd := exec.Command("dunstify", *errMsg, "-u", "critical")
		cmd.Run()
		os.Exit(1)
	}

	apiClient := api.Client{
		Domain:      authConfig.Domain,
		Email:       authConfig.Email,
		UserId:      authConfig.UserId,
		AdminApiKey: authConfig.AdminApiKey,
		ApiKey:      authConfig.ApiKey,
	}

	projects := apiClient.FetchActiveProjectsArr()

	customers := make(map[int]api.Customer)
	for _, project := range projects {
		customers[int(project.Customer.Id)] = project.Customer
	}

	activities := apiClient.FetchActivitiesArr("2023-04-09", "2023-04-09")

	fmt.Println("CUSTOMERS")
	fmt.Println(customers)
	fmt.Println("PROJECTS")
	fmt.Println(projects)
	fmt.Println("ACTIVITIES")
	fmt.Println(activities)

	// apiClient.CreateActivity(945824843, 12188024, "2023-04-09", 0, "TEST")

	// activity := apiClient.FetchActivity(1030703750)
	// fmt.Println(activity)

	// performanceReport := apiClient.FetchAnnualyVariationUntilToday()
	// fmt.Println(performanceReport)

	// apiClient.ControlActivityTimer(1030704248, "stop")

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
