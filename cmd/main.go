// main should be a small function that imports from /pkg and /internal
package main

import (
	"fmt"
	"os/user"
	"strconv"
	"time"

	"github.com/rodolfoksveiga/moco-menu/internal/config"
	"github.com/rodolfoksveiga/moco-menu/internal/templates"
	"github.com/rodolfoksveiga/moco-menu/pkg/api"
	"github.com/rodolfoksveiga/moco-menu/pkg/menu"
)

func main() {
	osUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	configDir := fmt.Sprintf("%s/.config/moco-menu", osUser.HomeDir)

	configInit := config.Init{ConfigPath: fmt.Sprintf("%s/config.json", configDir)}

	config := configInit.Load()

	templatesInit := templates.Init{TemplatePath: fmt.Sprintf("%s/templates.json", configDir)}
	templates := templatesInit.Load()
	fmt.Println("TEMPLATES")
	fmt.Println(templates)

	apiClient := api.Client{
		Domain:      config.Domain,
		Email:       config.Email,
		UserId:      config.UserId,
		AdminApiKey: config.AdminApiKey,
		ApiKey:      config.ApiKey,
	}

	projects := apiClient.FetchActiveProjectsArr()

	customers := make(map[int]api.Customer)
	for _, project := range projects {
		customers[int(project.Customer.Id)] = project.Customer
	}

	activities := apiClient.FetchActivitiesArr("2023-04-09", "2023-04-09")

	/* fmt.Println("CUSTOMERS")
	fmt.Println(customers)
	fmt.Println("PROJECTS")
	fmt.Println(projects)
	fmt.Println("ACTIVITIES")
	fmt.Println(activities) */

	mainMenuOptions := []string{"Control activity timer", "Create activity", "Edit activity", "Delete activity"}
	selected := menu.RunMenu("Moco menu:", mainMenuOptions)

	switch selected {
	case mainMenuOptions[0]:
		timerMenuOptions := []string{"Start", "Stop"}
		selected := menu.RunMenu("Timer control:", timerMenuOptions)

		switch selected {
		case timerMenuOptions[0]:
			var activitiesMenuOptions []string
			activitiesMenuMap := make(map[string]api.Activity)
			for index, activity := range activities {
				menuValue := fmt.Sprintf(
					"%s. %s - %s - %s - %s h - %s",
					strconv.Itoa(index),
					activity.Customer.Name,
					activity.Project.Name,
					activity.Task.Name,
					fmt.Sprintf("%.2f", activity.Hours),
					activity.Description,
				)
				activitiesMenuOptions = append(activitiesMenuOptions, menuValue)
				activitiesMenuMap[menuValue] = activity
			}
			selected := menu.RunMenu("Activity:", activitiesMenuOptions)

			apiClient.ControlActivityTimer(activitiesMenuMap[selected].Id, "start")
		case timerMenuOptions[1]:
			menu.RunConfirmationMenu()

			today := time.Now().Format("2006-01-02")

			activities := apiClient.FetchActivitiesArr(today, today)

			activityId := api.FindRunningActivityId(activities)

			if activityId != nil {
				apiClient.ControlActivityTimer(*activityId, "stop")
				menu.NotifySuccess("Timer stopped.")
			}

			menu.Exit("Timer was already not running.")
		default:
			menu.Exit("Invalid menu option.")
		}
	default:
		menu.Exit("Invalid menu option.")
	}
}
