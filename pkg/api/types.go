package api

type User struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
}

type ActivityProject struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Billable bool   `json:"billable"`
}

type Task struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Billable bool   `json:"billable"`
}

type Pattern struct {
	Am []float64 `json:"am"`
	Pm []float64 `json:"pm"`
}

type Employment struct {
	Id                int64   `json:"id"`
	WeeklyTargetHours float64 `json:"weekly_target_hours"`
	Pattern           Pattern `json:"pattern"`
	From              string  `json:"from"`
	To                string  `json:"to"`
	User              User    `json:"user"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

type Activity struct {
	Id             int64           `json:"id"`
	Date           string          `json:"date"`
	Hours          float64         `json:"hours"`
	Description    string          `json:"description,omitempty"`
	Billed         bool            `json:"billed"`
	Billable       bool            `json:"billable"`
	Project        ActivityProject `json:"project"`
	Task           Task            `json:"task"`
	Customer       Customer        `json:"customer"`
	TimerStartedAt string          `json:"timer_started_at"`
}

type ProjectTask struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Active   bool   `json:"active"`
	Billable bool   `json:"billable"`
}

type Project struct {
	Id       int64         `json:"id"`
	Name     string        `json:"name"`
	Active   bool          `json:"active"`
	Customer Customer      `json:"customer"`
	Tasks    []ProjectTask `json:"tasks"`
}

type Customer struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type AnnualReport struct {
	Year                int32   `json:"year"`
	EmploymentHours     float64 `json:"employment_hours"`
	TargetHours         float64 `json:"target_hours"`
	HoursTrackedTotal   float64 `json:"hours_tracked_total"`
	Variation           float64 `json:"variation"`
	VariationUntilToday float64 `json:"variation_until_today"`
}

type MonthlyReport struct {
	Year              int32   `json:"year"`
	Month             int32   `json:"month"`
	TargetHours       float64 `json:"target_hours"`
	HoursTrackedTotal float64 `json:"hours_tracked_total"`
	Variation         float64 `json:"variation"`
}

type PerformanceReport struct {
	Annually AnnualReport    `json:"annually"`
	Monthly  []MonthlyReport `json:"monthly"`
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

type ControlActivityTimer struct {
	Control    string `json:"control"`
	ActivityId int64  `json:"activity_id"`
}

type DeleteActivity struct {
	ActivityId int64 `json:"acitivity_id"`
}
