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

func (apiClient Client) FetUserId() *int64 {
	url := fmt.Sprintf("https://%s.mocoapp.com/api/v1/users", apiClient.Domain)
	authHeader := fmt.Sprintf("Token token=%s", apiClient.AdminApiKey)

	req, err := http.NewRequest("GET", url, nil)
	menu.ExitOnError(err, "Failed to create request.")

	prepareHeaders(req, authHeader)

	client := &http.Client{}

	resp, err := client.Do(req)
	menu.ExitOnError(err, "Request to fetch user id failed.")
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	menu.ExitOnError(err, "Failed to read response body.")

	var users []User
	json.Unmarshal(bodyBytes, &users)
	user := findUserByEmail(users, apiClient.Email)

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

func (apiClient Client) FetchActivities(from string, to string) []Activity {
	url := fmt.Sprintf(
		"https://%s.mocoapp.com/api/v1/activities?user_id=%s&from=%s&to=%s",
		apiClient.Domain,
		strconv.FormatInt(apiClient.UserId, 10),
		from,
		to,
	)
	authHeader := fmt.Sprintf("Token token=%s", apiClient.ApiKey)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	menu.ExitOnError(err, "Failed to create request.")

	prepareHeaders(req, authHeader)

	resp, err := client.Do(req)
	menu.ExitOnError(err, "Request to create activity failed.")
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	menu.ExitOnError(err, "Failed to read response body.")

	var activities []Activity
	json.Unmarshal(bodyBytes, &activities)

	return activities
}

func (apiClient Client) CreateActivity(
	projectId int64,
	taskId int64,
	date string,
	hours float64,
	description string,
) {
	url := fmt.Sprintf(
		"https://%s.mocoapp.com/api/v1/activities",
		apiClient.Domain,
	)
	authHeader := fmt.Sprintf("Token token=%s", apiClient.ApiKey)

	client := &http.Client{}

	newActivity := CreateActivity{
		ProjectId:   projectId,
		TaskId:      taskId,
		Date:        date,
		Hours:       hours,
		Description: description,
	}
	newActivityStr, err := json.Marshal(newActivity)
	menu.ExitOnError(err, "Failed to marshal request body.")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(newActivityStr))
	menu.ExitOnError(err, "Failed to create request.")

	prepareHeaders(req, authHeader)

	resp, err := client.Do(req)
	menu.ExitOnError(err, "Request to create activity failed.")
	defer resp.Body.Close()
}

func (apiClient Client) UpdateActivity(
	activityId int64,
	projectId int64,
	taskId int64,
	date string,
	hours float64,
	description string,
) {
	url := fmt.Sprintf(
		"https://%s.mocoapp.com/api/v1/activities/%s",
		apiClient.Domain,
		strconv.FormatInt(activityId, 10),
	)
	authHeader := fmt.Sprintf("Token token=%s", apiClient.ApiKey)

	client := &http.Client{}

	newActivity := UpdateActivity{
		ActivityId:  activityId,
		ProjectId:   projectId,
		TaskId:      taskId,
		Date:        date,
		Hours:       hours,
		Description: description,
	}
	newActivityStr, err := json.Marshal(newActivity)
	menu.ExitOnError(err, "Failed to marshal request body.")

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(newActivityStr))
	menu.ExitOnError(err, "Failed to create request.")

	prepareHeaders(req, authHeader)

	resp, err := client.Do(req)
	menu.ExitOnError(err, "Request to update activity failed.")
	defer resp.Body.Close()
}

func (apiClient Client) FetchAssignedProjects() []Project {
	url := fmt.Sprintf(
		"https://%s.mocoapp.com/api/v1/projects/assigned",
		apiClient.Domain,
	)
	authHeader := fmt.Sprintf("Token token=%s", apiClient.ApiKey)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	menu.ExitOnError(err, "Failed to create request.")

	prepareHeaders(req, authHeader)

	resp, err := client.Do(req)
	menu.ExitOnError(err, "Request to fetch assined projects failed.")
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	menu.ExitOnError(err, "Failed to read response body.")

	var projects []Project
	json.Unmarshal(bodyBytes, &projects)

	return projects
}

func (apiClient Client) FetchTasksByProjectId(projectId int64) []Task {
	url := fmt.Sprintf(
		"https://%s.mocoapp.com/api/v1/projects/%s/tasks",
		apiClient.Domain,
		strconv.FormatInt(projectId, 10),
	)
	authHeader := fmt.Sprintf("Token token=%s", apiClient.ApiKey)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	menu.ExitOnError(err, "Failed to create request.")

	prepareHeaders(req, authHeader)

	resp, err := client.Do(req)
	menu.ExitOnError(err, "Request to fetch tasks by project id failed.")
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	menu.ExitOnError(err, "Failed to read response body.")

	var tasks []Task
	json.Unmarshal(bodyBytes, &tasks)

	return tasks
}
