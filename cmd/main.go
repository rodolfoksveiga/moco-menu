// main should be a small function that imports from /pkg and /internal
package main

import (
	"fmt"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rodolfoksveiga/moco-menu/internal/config"
	"github.com/rodolfoksveiga/moco-menu/internal/templates"
	"github.com/rodolfoksveiga/moco-menu/pkg/api"
	"github.com/rodolfoksveiga/moco-menu/pkg/menu"
)

func prepareDateMenuOptions() []string {
	today := time.Now()
	dateMenuOptions := []string{
		fmt.Sprintf(
			"0. Today (%s)",
			today.Format("2006-01-02"),
		),
		fmt.Sprintf(
			"1. Yesterday (%s)",
			today.AddDate(0, 0, -1).Format("2006-01-02"),
		),
		fmt.Sprintf(
			"2. Two days ago (%s)",
			today.AddDate(0, 0, -2).Format("2006-01-02"),
		),
	}

	return dateMenuOptions
}

func prepareEndDateMenuOptions(startDateStr string) []string {
	startDate := menu.ValidateDate(startDateStr)
	endDateMenuOptions := []string{
		fmt.Sprintf(
			"0. Last input (%s)",
			startDate.Format("2006-01-02"),
		),
		fmt.Sprintf(
			"1. One day after (%s)",
			startDate.AddDate(0, 0, 1).Format("2006-01-02"),
		),
		fmt.Sprintf(
			"2. Two days after (%s)",
			startDate.AddDate(0, 0, 2).Format("2006-01-02"),
		),
	}

	return endDateMenuOptions
}

func prepareMenuOptions[T interface {
	api.Customer | api.Project | api.Task | templates.Template
	GetName() string
}](items []T) ([]string, map[string]T) {
	menuOptions := []string{}
	menuMap := make(map[string]T)
	for index, item := range items {
		nameOption := fmt.Sprintf("%02d. %s", index, item.GetName())
		if len(menuOptions) > 10 {
			nameOption = fmt.Sprintf("%s. %s", strconv.Itoa(index), item.GetName())
		}
		menuOptions = append(menuOptions, nameOption)
		menuMap[nameOption] = item
	}
	sort.Strings(menuOptions)
	return menuOptions, menuMap
}

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

	mainMenuOptions := menu.IndexMenuOptions([]string{
		"Control activity timer",
		"Create activity",
		"Update activity",
		"Delete activity",
	})
	selected := menu.RunMenu("Moco menu:", mainMenuOptions)

	switch selected {
	case mainMenuOptions[0]:
		timerMenuOptions := menu.IndexMenuOptions([]string{"Start", "Stop"})
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
		createActivityMenuOptions := menu.IndexMenuOptions(
			[]string{"Use template", "From scratch"},
		)
		selected := menu.RunMenu("Create activity:", createActivityMenuOptions)

		switch selected {
		case createActivityMenuOptions[0]:
			templatesInit := templates.Init{
				TemplatePath: fmt.Sprintf("%s/templates.json", configDir),
			}
			activityTemplates := templatesInit.Load()

			templateMenuOptions, templateMenuMap := prepareMenuOptions(*activityTemplates)
			selected := menu.RunMenu("Template:", templateMenuOptions)
			newActivityTemplate := templateMenuMap[selected]

			if newActivityTemplate.Duration == nil {
				selectedDuration := menu.RunMenu("Duration:", []string{})
				newDuration := menu.ValidateDuration(selectedDuration)
				newActivityTemplate.Duration = &newDuration
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

			customers := api.FilterUniqueCustomers(projects)
			customerMenuOptions, customerMenuMap := prepareMenuOptions(customers)
			selectedCustomer := menu.RunMenu("Costumer:", customerMenuOptions)
			newCustomerId := customerMenuMap[selectedCustomer].Id

			filteredProjects := api.FilterProjectsByCustomerId(projects, newCustomerId)
			projectMenuOptions, projectMenuMap := prepareMenuOptions(filteredProjects)
			selectedProject := menu.RunMenu("Project:", projectMenuOptions)
			newProjectId := projectMenuMap[selectedProject].Id

			tasks := apiClient.FetchProjectTasksArr(newProjectId)
			taskMenuOptions, taskMenuMap := prepareMenuOptions(tasks)
			selectedTask := menu.RunMenu("Task:", taskMenuOptions)
			newTaskId := taskMenuMap[selectedTask].Id

			dateMenuOptions := prepareDateMenuOptions()
			selectedDate := menu.RunMenu("Date:", dateMenuOptions)
			newDate := menu.ValidateDateOutput(selectedDate)

			selectedDuration := menu.RunMenu("Duration:", []string{})
			newDuration := menu.ValidateDuration(selectedDuration)

			selectedDescription := menu.RunMenu("Description:", []string{})
			newDescription := selectedDescription

			apiClient.CreateActivity(
				newProjectId,
				newTaskId,
				newDate,
				newDuration,
				newDescription,
			)
		default:
			menu.Exit("Invalid menu option.")
		}
	case mainMenuOptions[2]:
		startDateMenuOptions := prepareDateMenuOptions()
		selectedStartDate := menu.RunMenu("Date:", startDateMenuOptions)
		startDate := menu.ValidateDateOutput(selectedStartDate)
		fmt.Println(startDate)

		endDateMenuOptions := prepareEndDateMenuOptions(startDate)
		selectedEndDate := menu.RunMenu("Date:", endDateMenuOptions)
		endDate := menu.ValidateDateOutput(selectedEndDate)
		fmt.Println(endDate)

		activities := apiClient.FetchActivitiesArr(startDate, endDate)

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
		selectedActivity := menu.RunMenu("Activity:", activitiesMenuOptions)
		activity := activitiesMenuMap[selectedActivity]

		fmt.Println("ACTIVITY")
		fmt.Println(activity)

		updateActivityMenuOptions := menu.IndexMenuOptions([]string{
			"Finish update",
			"Description",
			"Duration",
			"Date",
			"Task",
			"Project",
		})
		for {
			selectedProperty := menu.RunMenu("Update property:", updateActivityMenuOptions)
			selectedPropertySlice := selectedProperty[3:]
			selectedPropertyDisplay := strings.ToLower(selectedPropertySlice)

			switch selectedProperty {
			case updateActivityMenuOptions[0]:
				apiClient.UpdateActivity(
					activity.Id,
					activity.Project.Id,
					activity.Task.Id,
					activity.Date,
					activity.Duration,
					activity.Description,
				)

				menu.NotifySuccess("Update success.")
			case updateActivityMenuOptions[1]:
				menuOptions := fmt.Sprintf(
					"Actual %s: (%s)",
					selectedPropertyDisplay,
					activity.Description,
				)
				newDescription := menu.RunMenu(
					fmt.Sprintf("Update %s:", selectedPropertyDisplay),
					[]string{menuOptions},
				)
				activity.Description = newDescription
			case updateActivityMenuOptions[2]:
				menuOptions := fmt.Sprintf(
					"Actual %s: (%s)",
					selectedPropertyDisplay,
					strconv.FormatFloat(activity.Duration, 'f', 2, 64),
				)

				selectedDuration := menu.RunMenu(
					fmt.Sprintf("Update %s:", selectedPropertyDisplay),
					[]string{menuOptions},
				)
				newDuration := menu.ValidateDuration(selectedDuration)
				activity.Duration = newDuration
			case updateActivityMenuOptions[3]:
				menuOptions := fmt.Sprintf(
					"Actual %s: (%s)",
					selectedPropertyDisplay,
					activity.Date,
				)
				selectedDate := menu.RunMenu(
					fmt.Sprintf("Update %s:", selectedPropertyDisplay),
					[]string{menuOptions},
				)
				newDateStr := menu.ValidateDateOutput(selectedDate)
				activity.Date = newDateStr
			case updateActivityMenuOptions[4]:
				tasks := apiClient.FetchProjectTasksArr(activity.Project.Id)

				// firstMenuOption := fmt.Sprintf("Actual %s: (%s)", selectedPropertyDisplay, activity.Project.Name)
				taskMenuOptions, taskMenuMap := prepareMenuOptions(tasks)
				// menuOptions = append([]string{firstMenuOption}, menuOptions...)
				selectedTask := menu.RunMenu(fmt.Sprintf("Update %s:", selectedPropertyDisplay), taskMenuOptions)
				newTaskId := taskMenuMap[selectedTask].Id
				fmt.Sprintln(strconv.FormatInt(newTaskId, 10))
				// activity.Project.Id = newProjectId
			case updateActivityMenuOptions[5]:
				projects := apiClient.FetchActiveProjectsArr()
				filteredProjects := api.FilterProjectsByCustomerId(projects, activity.Customer.Id)

				// firstMenuOption := fmt.Sprintf("Actual %s: (%s)", selectedPropertyDisplay, activity.Project.Name)
				projectMenuOptions, projectMenuMap := prepareMenuOptions(filteredProjects)
				// menuOptions = append([]string{firstMenuOption}, menuOptions...)
				selectedProject := menu.RunMenu(fmt.Sprintf("Update %s:", selectedPropertyDisplay), projectMenuOptions)
				newProjectId := projectMenuMap[selectedProject].Id
				fmt.Sprintln(strconv.FormatInt(newProjectId, 10))
				activity.Project.Id = newProjectId
			default:
				menu.Exit("Invalid menu option.")
			}
		}

	case mainMenuOptions[3]:
		startDateMenuOptions := prepareDateMenuOptions()
		selectedStartDate := menu.RunMenu("Date:", startDateMenuOptions)
		startDate := menu.ValidateDateOutput(selectedStartDate)
		fmt.Println(startDate)

		endDateMenuOptions := prepareEndDateMenuOptions(startDate)
		selectedEndDate := menu.RunMenu("Date:", endDateMenuOptions)
		endDate := menu.ValidateDateOutput(selectedEndDate)
		fmt.Println(endDate)

		activities := apiClient.FetchActivitiesArr(startDate, endDate)

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
		selectedActivity := menu.RunMenu("Activity:", activitiesMenuOptions)
		activity := activitiesMenuMap[selectedActivity]

		menu.RunConfirmationMenu()

		apiClient.DeleteActivity(activity.Id)
	default:
		menu.Exit("Invalid menu option.")
	}
}
