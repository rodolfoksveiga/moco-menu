package config

type AuthConfig struct {
	Domain      string `json:"domain"`
	Email       string `json:"email"`
	UserId      int64  `json:"userId"`
	AdminApiKey string `json:"adminApiKey"`
	ApiKey      string `json:"apiKey"`
	// UserId      *int64 `json:"userId,omitempty"`
}
