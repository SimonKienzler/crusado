package config

import (
	"log"
	"os"
)

const (
	OrganizationURLEnvVarKey = "CRUSADO_AZURE_ORG_URL"
	AzurePATEnvVarKey        = "CRUSADO_AZURE_PAT"
	ProjectNameEnvVarKey     = "CRUSADO_AZURE_PROJECT_NAME"
	TemplatesDirEnvVarKey    = "CRUSADO_TEMPLATES_DIR"
)

type Crusado struct {
	OrganizationURL     string
	PersonalAccessToken string
	ProjectName         string
	TemplatesDirectory  string
}

func GetConfigOrDie() Crusado {
	cfg := Crusado{}
	incomplete := false

	if organizationURL, exists := os.LookupEnv(OrganizationURLEnvVarKey); exists {
		cfg.OrganizationURL = organizationURL
	} else {
		incomplete = true
		log.Printf("Required environment variable %s is not set", OrganizationURLEnvVarKey)
	}
	if personalAccessToken, exists := os.LookupEnv(AzurePATEnvVarKey); exists {
		cfg.PersonalAccessToken = personalAccessToken
	} else {
		incomplete = true
		log.Printf("Required environment variable %s is not set", AzurePATEnvVarKey)
	}
	if projectName, exists := os.LookupEnv(ProjectNameEnvVarKey); exists {
		cfg.ProjectName = projectName
	} else {
		incomplete = true
		log.Printf("Required environment variable %s is not set", ProjectNameEnvVarKey)
	}
	if templatesDirectory, exists := os.LookupEnv(TemplatesDirEnvVarKey); exists {
		cfg.TemplatesDirectory = templatesDirectory
	} else {
		incomplete = true
		log.Printf("Required environment variable %s is not set", TemplatesDirEnvVarKey)
	}

	// TODO check if TemplatesDirectory is actually a directory

	if incomplete {
		os.Exit(1)
	}

	return cfg
}
