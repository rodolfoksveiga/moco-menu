// main should be a small function that imports from /pkg and /internal
package main

import (
	"fmt"
	"os/user"
	"regexp"
	"sort"
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

	apiClient := api.Client{
		Domain:      config.Domain,
		Email:       config.Email,
		UserId:      config.UserId,
		AdminApiKey: config.AdminApiKey,
		ApiKey:      config.ApiKey,
	}

	mainMenuOptions := []string{
		"Control activity timer",
		"Create activity",
		"Edit activity",
		"Delete activity",
	}
	selected := menu.RunMenu("Moco menu:", mainMenuOptions)

	switch selected {
	case mainMenuOptions[0]:
		timerMenuOptions := []string{"Start", "Stop"}
		selected := menu.RunMenu("Timer control:", timerMenuOptions)

		switch selected {
		case timerMenuOptions[0]:
			today := time.Now().Format("2006-01-02")
			activities := apiClient.FetchActivitiesArr(today, today)

			var activitiesMenuOptions []string
			activitiesMenuMap := make(map[string]api.Activity)
			for index, activity := range activities {
				menuValue := fmt.Sprintf(
					"%s. %s - %s - %s - %s h - %s",
					strconv.Itoa(index),
					activity.Customer.Name,
					activity.Project.Name,
					activity.Task.Name,
					fmt.Sprintf("%.2f", activity.Duration),
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
	case mainMenuOptions[1]:
		createActivityMenuOptions := []string{"Use template", "From scratch"}
		selected := menu.RunMenu("Create activity:", createActivityMenuOptions)

		switch selected {
		case createActivityMenuOptions[0]:
			templatesInit := templates.Init{
				TemplatePath: fmt.Sprintf("%s/templates.json", configDir),
			}
			activityTemplate := templatesInit.Load()

			chooseTemplateMenuOptions := []string{}

			chooseTemplateMenuMap := make(map[string]templates.Template)
			for _, template := range *activityTemplate {
				chooseTemplateMenuOptions = append(
					chooseTemplateMenuOptions,
					template.Name,
				)
				chooseTemplateMenuMap[template.Name] = template
			}
			selected := menu.RunMenu("Template:", chooseTemplateMenuOptions)
			newActivityTemplate := chooseTemplateMenuMap[selected]

			if newActivityTemplate.Duration == nil {
				selected := menu.RunMenu("Duration:", []string{})
				selectedFloat, err := strconv.ParseFloat(selected, 64)
				menu.ExitOnError(err, "Failed to parse duration.")
				newActivityTemplate.Duration = &selectedFloat
			}

			if newActivityTemplate.Description == nil {
				selected := menu.RunMenu("Description:", []string{})
				newActivityTemplate.Description = &selected
			}

			apiClient.CreateActivity(
				newActivityTemplate.ProjectId,
				newActivityTemplate.TaskId,
				time.Now().Format("2006-01-02"),
				*newActivityTemplate.Duration,
				*newActivityTemplate.Description,
			)

		case createActivityMenuOptions[1]:
			projects := apiClient.FetchActiveProjectsArr()

			customers := make(map[int64]api.Customer)
			for _, project := range projects {
				customers[project.Customer.Id] = project.Customer
			}

			chooseCustomerMenuOptions := []string{}
			chooseCustomerMenuMap := make(map[string]api.Customer)
			for _, customer := range customers {
				chooseCustomerMenuOptions = append(
					chooseCustomerMenuOptions,
					customer.Name,
				)
				chooseCustomerMenuMap[customer.Name] = customer
			}
			sort.Strings(chooseCustomerMenuOptions)
			selectedCustomer := menu.RunMenu("Costumer:", chooseCustomerMenuOptions)
			newCustomerId := chooseCustomerMenuMap[selectedCustomer].Id

			filteredProjects := api.FilterProjectsByCustomerId(projects, newCustomerId)

			chooseProjectMenuOptions := []string{}
			chooseProjectMenuMap := make(map[string]api.Project)
			for _, project := range filteredProjects {
				chooseProjectMenuOptions = append(
					chooseProjectMenuOptions,
					project.Name,
				)
				chooseProjectMenuMap[project.Name] = project
			}
			sort.Strings(chooseProjectMenuOptions)
			selectedProject := menu.RunMenu("Project:", chooseProjectMenuOptions)
			newProjectId := chooseProjectMenuMap[selectedProject].Id

			tasks := apiClient.FetchProjectTasksArr(newProjectId)
			chooseTaskMenuOptions := []string{}
			chooseTaskMenuMap := make(map[string]api.Task)
			for _, task := range tasks {
				chooseTaskMenuOptions = append(
					chooseTaskMenuOptions,
					task.Name,
				)
				chooseTaskMenuMap[task.Name] = task
			}
			sort.Strings(chooseTaskMenuOptions)
			selectedTask := menu.RunMenu("Task:", chooseTaskMenuOptions)
			newTaskId := chooseTaskMenuMap[selectedTask].Id

			today := time.Now()
			chooseDateMenuOptions := []string{
				fmt.Sprintf(
					"Today (%s)",
					today.Format("2006-01-02"),
				),
				fmt.Sprintf(
					"Yesterday (%s)",
					today.AddDate(0, 0, -1).Format("2006-01-02"),
				),
				fmt.Sprintf(
					"Before yesterday (%s)",
					today.AddDate(0, 0, -2).Format("2006-01-02"),
				),
			}
			newDateOutput := menu.RunMenu("Date:", chooseDateMenuOptions)
			re := regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}")
			newDate := re.FindString(newDateOutput)
			_, err := time.Parse("2006-01-02", newDate)
			menu.ExitOnError(err, "Wrong date format.")

			selectedDuration := menu.RunMenu("Duration:", []string{})
			newDuration, err := strconv.ParseFloat(selectedDuration, 64)
			menu.ExitOnError(err, "Wrong duration format.")

			selectedDescription := menu.RunMenu("Description:", []string{})
			newDescription := selectedDescription

			apiClient.CreateActivity(
				newProjectId,
				newTaskId,
				newDate,
				newDuration,
				newDescription,
			)
		case createActivityMenuOptions[2]:
			menu.Exit("Work in progress.")
		case createActivityMenuOptions[3]:
			menu.Exit("Work in progress.")
		default:
			menu.Exit("Invalid menu option.")
		}

	default:
		menu.Exit("Invalid menu option.")
	}
}
