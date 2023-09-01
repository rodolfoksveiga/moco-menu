package templates

type Template struct {
	Name        string   `json:"name"`
	ProjectId   int64    `json:"projectId"`
	TaskId      int64    `json:"taskId"`
	Duration    *float64 `json:"duration,omitempty"`
	Description *string  `json:"description,omitempty"`
}

func (t Template) GetName() string { return t.Name }
