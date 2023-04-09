package templates

type Template struct {
	Name        string  `json:"name"`
	Customer    string  `json:"customer,omitempty"`
	Project     string  `json:"project,omitempty"`
	Task        string  `json:"task,omitempty"`
	Duration    float64 `json:"duration,omitempty"`
	Description string  `json:"description,omitempty"`
}
