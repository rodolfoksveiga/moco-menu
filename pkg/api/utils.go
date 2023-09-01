package api

import (
	"net/http"
)

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

func FilterProjectsByCustomerId(projects []Project, customerId int64) []Project {
	var filteredProjects []Project
	for _, project := range projects {
		if project.Customer.Id == customerId {
			filteredProjects = append(filteredProjects, project)
		}
	}

	return filteredProjects
}

func FilterUniqueCustomers(projects []Project) []Customer {
	var customers []Customer
	for _, project := range projects {
		hasCustomer := false
		for _, customer := range customers {
			if customer == project.Customer {
				hasCustomer = true
				break
			}
		}

		if !hasCustomer {
			customers = append(customers, project.Customer)
		}
	}

	return customers
}
