package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/rodolfoksveiga/moco-menu/pkg/menu"
)

type Client struct {
	Domain      string
	Email       string
	UserId      int64
	AdminApiKey string
	ApiKey      string
}

func (client Client) FetchUserId() *int64 {
	url := fmt.Sprintf("https://%s.mocoapp.com/api/v1/users", client.Domain)
	authHeader := fmt.Sprintf("Token token=%s", client.AdminApiKey)

	req, err := http.NewRequest("GET", url, nil)
	menu.ExitOnError(err, "Failed to create request.")

	prepareHeaders(req, authHeader)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	menu.ExitOnError(err, "Request to fetch user id failed.")
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	menu.ExitOnError(err, "Failed to read response body.")

	var users []User
	json.Unmarshal(bodyBytes, &users)
	user := findUserByEmail(users, client.Email)

	if user == nil {
		cmd := exec.Command(
			"dunstify",
			"Could not find user. Review your config.",
			"-u",
			"critical",
		)
		cmd.Run()
		os.Exit(1)
	}

	return &user.Id
}

func (client Client) FetchAnnualyVariationUntilToday() float64 {
	url := fmt.Sprintf(
		"https://%s.mocoapp.com/api/v1/users/%s/performance_report",
		client.Domain,
		strconv.FormatInt(client.UserId, 10),
	)
	authHeader := fmt.Sprintf("Token token=%s", client.AdminApiKey)

	httpClient := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	menu.ExitOnError(err, "Failed to create request.")

	prepareHeaders(req, authHeader)

	resp, err := httpClient.Do(req)
	menu.ExitOnError(err, "Request to fetch annualy variation until today failed.")
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	menu.ExitOnError(err, "Failed to read response body.")

	var performanceReport PerformanceReport
	json.Unmarshal(bodyBytes, &performanceReport)

	return performanceReport.Annually.VariationUntilToday
}

func (client Client) FetchActiveProjectsArr() []Project {
	url := fmt.Sprintf(
		"https://%s.mocoapp.com/api/v1/projects/assigned?active=true",
		client.Domain,
	)
	authHeader := fmt.Sprintf("Token token=%s", client.ApiKey)

	httpClient := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	menu.ExitOnError(err, "Failed to create request.")

	prepareHeaders(req, authHeader)

	resp, err := httpClient.Do(req)
	menu.ExitOnError(err, "Request to fetch projects failed.")
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	menu.ExitOnError(err, "Failed to read response body.")

	var projects []Project
	json.Unmarshal(bodyBytes, &projects)

	return projects
}

func (client Client) FetchProjectTasksArr(projectId int64) []Task {
	url := fmt.Sprintf(
		"https://%s.mocoapp.com/api/v1/projects/%s/tasks",
		client.Domain,
		strconv.FormatInt(projectId, 10),
	)
	authHeader := fmt.Sprintf("Token token=%s", client.ApiKey)

	httpClient := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	menu.ExitOnError(err, "Failed to create request.")

	prepareHeaders(req, authHeader)

	resp, err := httpClient.Do(req)
	menu.ExitOnError(err, "Request to fetch project tasks id failed.")
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	menu.ExitOnError(err, "Failed to read response body.")

	var tasks []Task
	json.Unmarshal(bodyBytes, &tasks)

	return tasks
}

func (client Client) FetchActivity(activityId int64) Activity {
	url := fmt.Sprintf(
		"https://%s.mocoapp.com/api/v1/activities/%s",
		client.Domain,
		strconv.FormatInt(activityId, 10),
	)
	authHeader := fmt.Sprintf("Token token=%s", client.ApiKey)

	httpClient := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	menu.ExitOnError(err, "Failed to create request.")

	prepareHeaders(req, authHeader)

	resp, err := httpClient.Do(req)
	menu.ExitOnError(err, "Request to fecth activity failed.")
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	menu.ExitOnError(err, "Failed to read response body.")

	var activity Activity
	json.Unmarshal(bodyBytes, &activity)

	return activity
}

func (client Client) FetchActivitiesArr(from string, to string) []Activity {
	url := fmt.Sprintf(
		"https://%s.mocoapp.com/api/v1/activities?user_id=%s&from=%s&to=%s",
		client.Domain,
		strconv.FormatInt(client.UserId, 10),
		from,
		to,
	)
	authHeader := fmt.Sprintf("Token token=%s", client.ApiKey)

	httpClient := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	menu.ExitOnError(err, "Failed to create request.")

	prepareHeaders(req, authHeader)

	resp, err := httpClient.Do(req)
	menu.ExitOnError(err, "Request to fetch activities failed.")
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	menu.ExitOnError(err, "Failed to read response body.")

	var activities []Activity
	json.Unmarshal(bodyBytes, &activities)

	return activities
}

func (client Client) CreateActivity(
	projectId int64,
	taskId int64,
	date string,
	duration float64,
	description string,
) {
	url := fmt.Sprintf(
		"https://%s.mocoapp.com/api/v1/activities",
		client.Domain,
	)
	authHeader := fmt.Sprintf("Token token=%s", client.ApiKey)

	httpClient := &http.Client{}

	newActivity := CreateActivity{
		ProjectId:   projectId,
		TaskId:      taskId,
		Date:        date,
		Duration:    duration,
		Description: description,
	}
	newActivityStr, err := json.Marshal(newActivity)
	menu.ExitOnError(err, "Failed to marshal request body.")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(newActivityStr))
	menu.ExitOnError(err, "Failed to create request.")

	prepareHeaders(req, authHeader)

	resp, err := httpClient.Do(req)
	menu.ExitOnError(err, "Request to create activity failed.")
	defer resp.Body.Close()
}

func (client Client) UpdateActivity(
	activityId int64,
	projectId int64,
	taskId int64,
	date string,
	duration float64,
	description string,
) {
	url := fmt.Sprintf(
		"https://%s.mocoapp.com/api/v1/activities/%s",
		client.Domain,
		strconv.FormatInt(activityId, 10),
	)
	authHeader := fmt.Sprintf("Token token=%s", client.ApiKey)

	httpClient := &http.Client{}

	newActivity := UpdateActivity{
		Id:          activityId,
		ProjectId:   projectId,
		TaskId:      taskId,
		Date:        date,
		Duration:    duration,
		Description: description,
	}
	newActivityStr, err := json.Marshal(newActivity)
	menu.ExitOnError(err, "Failed to marshal request body.")

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(newActivityStr))
	menu.ExitOnError(err, "Failed to create request.")

	prepareHeaders(req, authHeader)

	resp, err := httpClient.Do(req)
	menu.ExitOnError(err, "Request to update activity failed.")
	defer resp.Body.Close()
}

func (client Client) DeleteActivity(
	activityId int64,
) {
	url := fmt.Sprintf(
		"https://%s.mocoapp.com/api/v1/activities/%s",
		client.Domain,
		strconv.FormatInt(activityId, 10),
	)
	authHeader := fmt.Sprintf("Token token=%s", client.ApiKey)

	httpClient := &http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	menu.ExitOnError(err, "Failed to create request.")

	prepareHeaders(req, authHeader)

	resp, err := httpClient.Do(req)
	menu.ExitOnError(err, "Request to delete activity failed.")
	defer resp.Body.Close()
}

func (client Client) ControlActivityTimer(
	activityId int64,
	control string,
) {
	if control == "start" || control == "stop" {
		url := fmt.Sprintf(
			"https://%s.mocoapp.com/api/v1/activities/%s/%s_timer",
			client.Domain,
			strconv.FormatInt(activityId, 10),
			control,
		)
		authHeader := fmt.Sprintf("Token token=%s", client.ApiKey)

		httpClient := &http.Client{}

		req, err := http.NewRequest("PATCH", url, nil)
		menu.ExitOnError(err, "Failed to create request.")

		prepareHeaders(req, authHeader)

		resp, err := httpClient.Do(req)
		menu.ExitOnError(
			err,
			fmt.Sprintf("Request to %s activity timer failed.", control),
		)
		defer resp.Body.Close()
	} else {
		menu.Exit("Timer control can't be different than \"start\" or \"stop\".")
	}
}
