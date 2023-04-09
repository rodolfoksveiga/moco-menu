package api

import "net/http"

func prepareHeaders(req *http.Request, authHeader string) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
}

func findUserByEmail(users []User, email string) *User {
	for i := 0; i < len(users); i++ {
		if users[i].Email == email {
			return &users[i]
		}
	}
	return nil
}

func FindRunningActivityId(activities []Activity) *int64 {
	var activityId int64
	for i := 0; i < len(activities); i++ {
		if activities[i].TimerStartedAt != nil {
			activityId = activities[i].Id
		}
	}

	return &activityId
}
