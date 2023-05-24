package config

import "os"

type CrusadoConfig struct {
	OrganizationUrl            string
	PersonalAccessToken        string
	ProjectName                string
	IterationPath              string
	UseIterationPathFromEnvVar bool
}

type ProjectConfig struct {
	Name          string
	IterationPath string
	AreaPath      string
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
