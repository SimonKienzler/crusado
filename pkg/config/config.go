package config

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type CrusadoConfig struct {
	OrganizationUrl            string
	PersonalAccessToken        string
	ProjectName                string
	IterationPath              string
	UseIterationPathFromEnvVar bool
}

func GetConfig(useIterationPathFromEnvVar bool) CrusadoConfig {
	organizationUrl := os.Getenv("AZURE_ORG_URL")
	personalAccessToken := os.Getenv("AZURE_PAT")
	projectName := os.Getenv("AZURE_PROJECT_NAME")
	currentIteration := os.Getenv("ITERATION_PATH")

	return CrusadoConfig{
		OrganizationUrl:            organizationUrl,
		PersonalAccessToken:        personalAccessToken,
		ProjectName:                projectName,
		IterationPath:              currentIteration,
		UseIterationPathFromEnvVar: useIterationPathFromEnvVar,
	}
}

type ProjectConfig struct {
	Name          string
	IterationPath string
	AreaPath      string
}

type Profile struct {
	Name      string              `yaml:"name"`
	Templates []UserStoryTemplate `yaml:"templates"`
}

type UserStoryTemplate struct {
	Name             string `yaml:"name" json:"name"`
	Description      string `yaml:"description" json:"description"`
	StoryTitle       string `yaml:"storyTitle" json:"storyTitle"`
	StoryDescription string `yaml:"storyDescription" json:"storyDescription"`
	Tasks            []Task `yaml:"tasks" json:"tasks"`
}

type Task struct {
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description" json:"description"`
}

func GetProfileFromFile(filepath string) (*Profile, error) {
	exampleProfile, err := ioutil.ReadFile("./example/profile.yaml")
	if err != nil {
		return nil, err
	}

	profile := Profile{}
	err = yaml.Unmarshal(exampleProfile, &profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}
