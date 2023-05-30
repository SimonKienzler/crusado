package config

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/thediveo/klo"
	yaml "gopkg.in/yaml.v3"
)

const (
	OrganizationUrlEnvVarKey     = "CRUSADO_AZURE_ORG_URL"
	PersonalAccessTokenEnvVarKey = "CRUSADO_AZURE_PAT"
	ProjectNameEnvVarKey         = "CRUSADO_AZURE_PROJECT_NAME"
	IterationPathEnvVarKey       = "CRUSADO_ITERATION_PATH"
	ProfileFilePathEnvVarKey     = "CRUSADO_PROFILE_FILE_PATH"
	UseCurrentIterationEnvVarKey = "CRUSADO_USE_CURRENT_ITERATION"
)

type CrusadoConfig struct {
	OrganizationUrl     string
	PersonalAccessToken string
	ProjectName         string
	IterationPath       string
	UseCurrentIteration bool
	ProfileFilePath     string
}

func GetConfig() CrusadoConfig {
	cfg := CrusadoConfig{}
	incomplete := false

	if organizationUrl, exists := os.LookupEnv(OrganizationUrlEnvVarKey); exists {
		cfg.OrganizationUrl = organizationUrl
	} else {
		incomplete = true
		log.Printf("Required environment variable %s is not set", OrganizationUrlEnvVarKey)
	}
	if personalAccessToken, exists := os.LookupEnv(PersonalAccessTokenEnvVarKey); exists {
		cfg.PersonalAccessToken = personalAccessToken
	} else {
		incomplete = true
		log.Printf("Required environment variable %s is not set", PersonalAccessTokenEnvVarKey)
	}
	if projectName, exists := os.LookupEnv(ProjectNameEnvVarKey); exists {
		cfg.ProjectName = projectName
	} else {
		incomplete = true
		log.Printf("Required environment variable %s is not set", ProjectNameEnvVarKey)
	}
	if iterationPath, exists := os.LookupEnv(IterationPathEnvVarKey); exists {
		cfg.IterationPath = iterationPath
	} else {
		incomplete = true
		log.Printf("Required environment variable %s is not set", IterationPathEnvVarKey)
	}
	if profileFilePath, exists := os.LookupEnv(ProfileFilePathEnvVarKey); exists {
		cfg.ProfileFilePath = profileFilePath
	} else {
		incomplete = true
		log.Printf("Required environment variable %s is not set", ProfileFilePathEnvVarKey)
	}
	if useCurrentIteration, exists := os.LookupEnv(UseCurrentIterationEnvVarKey); exists {
		cfg.UseCurrentIteration = useCurrentIteration == "true"
	} else {
		incomplete = true
		log.Printf("Required environment variable %s is not set", UseCurrentIterationEnvVarKey)
	}

	if incomplete {
		os.Exit(1)
	}

	return cfg
}

type ProjectConfig struct {
	Name          string
	IterationPath string
	AreaPath      string
}

type TemplateList struct {
	Templates []Template `yaml:"templates" json:"templates"`
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

type Template struct {
	// Name is the unique name of the template, used in commands
	Name string `yaml:"name" json:"name"`

	// Summary provides a short synopsis for the template
	Summary string `yaml:"summary" json:"summary"`

	// Type identifies the template as one of [User Story, Bug]
	Type TemplateType `yaml:"type" json:"type"`

	// Title is the resulting title of the work item in Azure DevOps
	Title string `yaml:"title" json:"title"`

	// Description contains the resulting text in the body of the work item
	Description string `yaml:"description" json:"description"`

	// Tasks is a slice of individual tasks that are part of the template
	Tasks []Task `yaml:"tasks" json:"tasks"`
}

type TemplateType string

const (
	TemplateTypeUserStory = TemplateType("UserStory")
	TemplateTypeBug       = TemplateType("Bug")
	TemplateTypeTask      = TemplateType("Task")
)

var AvailableTemplateTypes = []TemplateType{
	TemplateTypeUserStory,
	TemplateTypeBug,
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

	templateList := TemplateList{}
	err = yaml.Unmarshal(exampleProfile, &templateList)
	if err != nil {
		return nil, err
	}

	return &templateList, nil
}

var TemplatePrinterSpecs = klo.Specs{
	DefaultColumnSpec: "NAME:{.Name},TYPE:{.Type},SUMMARY:{.Summary}",
	WideColumnSpec:    "NAME:{.Name},TYPE:{.Type},SUMMARY:{.Summary},TITLE:{.Title},TASKS:{.Tasks[*].Title}",
}
