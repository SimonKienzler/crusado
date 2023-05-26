package config

import (
	"io/ioutil"
	"os"

	"github.com/thediveo/klo"
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

type TemplateList struct {
	Templates []UserStoryTemplate `yaml:"templates" json:"templates"`
}

type Profile struct {
	Name              string `yaml:"name" json:"name"`
	FilePath          string `yaml:"filePath" json:"filePath"`
	Valid             bool   `yaml:"valid" json:"valid"`
	NumberOfTemplates int    `yaml:"numberOfTemplates" json:"numberOfTemplates"`
}

type ProfileConfig struct {
	Name     string `yaml:"name" json:"name"`
	FilePath string `yaml:"filePath" json:"filePath"`
}

var ProfilePrinterSpecs = klo.Specs{
	DefaultColumnSpec: "NAME:{.Name},FILEPATH:{.FilePath},VALID:{.Valid}",
	WideColumnSpec:    "NAME:{.Name},FILEPATH:{.FilePath},VALID:{.Valid},TEMPLATES:{.NumberOfTemplates}",
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

func GetTemplateListFromFile(filepath string) (*TemplateList, error) {
	exampleProfile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	profile := TemplateList{}
	err = yaml.Unmarshal(exampleProfile, &profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

var TemplatePrinterSpecs = klo.Specs{
	DefaultColumnSpec: "NAME:{.Name},DESCRIPTION:{.Description}",
	WideColumnSpec:    "NAME:{.Name},DESCRIPTION:{.Description},STORY TITLE:{.StoryTitle},TASKS:{.Tasks[*].Title}",
}
