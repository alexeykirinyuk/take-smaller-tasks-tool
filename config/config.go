package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type SMTPConfiguration struct {
	UserName string
	Password string
	Domain   string
	Port     int
}

type Configuration struct {
	Login  string
	Token  string
	Domain string

	Large                   []string
	Small                   []string
	IssuesByEstimation      string
	AllIssuesWithEstimation string

	EmailNotificationsEnabled bool
	SMTP                      SMTPConfiguration
}

func Get() (config Configuration, err error) {
	file, err := os.Open("config.json")

	if err != nil {
		err = fmt.Errorf("error when open configuration file: %s", err)
		return
	}

	defer func() {
		err = file.Close()
	}()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	if err != nil {
		err = fmt.Errorf("error when parsing configurations: %s", err)
		return
	}

	return
}
