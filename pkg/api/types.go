package api

type User struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
}

type Customer struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Project struct {
	Id       int64    `json:"id"`
	Name     string   `json:"name"`
	Customer Customer `json:"customer"`
}

type ActivityProject struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Task struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Activity struct {
	Id             int64           `json:"id"`
	Date           string          `json:"date"`
	Hours          float64         `json:"hours"`
	Description    string          `json:"description,omitempty"`
	Project        ActivityProject `json:"project"`
	Task           Task            `json:"task"`
	Customer       Customer        `json:"customer"`
	TimerStartedAt *string         `json:"timer_started_at"`
}

type AnnualReport struct {
	VariationUntilToday float64 `json:"variation_until_today"`
}

type PerformanceReport struct {
	Annually AnnualReport `json:"annually"`
}

type CreateActivity struct {
	ProjectId   int64   `json:"project_id"`
	TaskId      int64   `json:"task_id"`
	Date        string  `json:"date"`
	Hours       float64 `json:"hours,omitempty"`
	Description string  `json:"description"`
}

type UpdateActivity struct {
	ActivityId  int64   `json:"activity_id"`
	ProjectId   int64   `json:"project_id"`
	TaskId      int64   `json:"task_id"`
	Date        string  `json:"date"`
	Hours       float64 `json:"hours"`
	Description string  `json:"description"`
}
