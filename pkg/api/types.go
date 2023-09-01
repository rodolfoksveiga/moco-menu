package api

type User struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
}

type Customer struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (t Customer) GetId() int64    { return t.Id }
func (t Customer) GetName() string { return t.Name }

type Project struct {
	Id       int64    `json:"id"`
	Name     string   `json:"name"`
	Customer Customer `json:"customer"`
}

func (t Project) GetId() int64    { return t.Id }
func (t Project) GetName() string { return t.Name }

type ActivityProject struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Task struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (t Task) GetId() int64    { return t.Id }
func (t Task) GetName() string { return t.Name }

func GetName[T interface{ GetName() string }](t T) string {
	return t.GetName()
}

type Activity struct {
	Id             int64           `json:"id"`
	Date           string          `json:"date"`
	Duration       float64         `json:"hours"`
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
	Duration    float64 `json:"hours,omitempty"`
	Description string  `json:"description"`
}

type UpdateActivity struct {
	Id          int64   `json:"activity_id"`
	ProjectId   int64   `json:"project_id"`
	TaskId      int64   `json:"task_id"`
	Date        string  `json:"date"`
	Duration    float64 `json:"hours"`
	Description string  `json:"description"`
}
