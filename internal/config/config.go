package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	ConfigFilePath string
}

func (config Config) Load() (*AuthConfig, *string) {
	configFile, err := os.Open(config.ConfigFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		errMsg := "Could not find config file."
		return nil, &errMsg
	}
	defer configFile.Close()

	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		fmt.Println("Error:", err)
		errMsg := "Could not read config file."
		return nil, &errMsg
	}

	var authConfig AuthConfig
	err = json.Unmarshal(configBytes, &authConfig)
	if err != nil {
		fmt.Println("Error:", err)
		errMsg := "Could not map config file."
		return nil, &errMsg
	}

	/* if authConfig.UserId == nil {
		apiClient := api.Client{Domain: authConfig.Domain, Email: authConfig.Email, AdminApiKey: authConfig.AdminApiKey, ApiKey: authConfig.ApiKey}
		userId := apiClient.FetUserId()
		authConfig.UserId = userId

		authConfigFile, err := json.MarshalIndent(authConfig, "", "  ")
		if err != nil {
			fmt.Println("Error:", err)
			errMsg := "Could not read config file."
			return nil, &errMsg
		}

		err = ioutil.WriteFile(config.ConfigFilePath, authConfigFile, 0644)
		if err != nil {
			fmt.Println("Error:", err)
			errMsg := "Could not write config file."
			return nil, &errMsg
		}
	} */

	return &authConfig, nil
}

/* func (config Config) login() (*AuthConfig, *string) {
	domain := menu.RunMenu("Domain:", []string{})
	email := menu.RunMenu("Email:", []string{})
	userId := menu.RunMenu("UserId:", []string{})
	adminApiKey := menu.RunMenu("Admin API Key:", []string{})
	apiKey := menu.RunMenu("API Key:", []string{})

	authConfig := AuthConfig{Domain: domain, Email: email, UserId: userId, AdminApiKey: adminApiKey, ApiKey: apiKey}
	authConfigFile, err := json.MarshalIndent(authConfig, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		errMsg := "Could not map config file."
		return nil, &errMsg
	}

	confirm := menu.RunMenu("Confirm action:", []string{"Yes", "No"})
	if strings.Contains(confirm, "Yes") {
		err = ioutil.WriteFile(config.ConfigFilePath, authConfigFile, 0644)
		if err != nil {
			fmt.Println("Error:", err)
			errMsg := "Could not write config file."
			return nil, &errMsg
		}

		return &authConfig, nil
	}

	errMsg := "User opted to not write the file."
	return nil, &errMsg
} */
