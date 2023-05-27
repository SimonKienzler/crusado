package profile

import (
	"log"
	"os"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/simonkienzler/crusado/pkg/validator"
	"github.com/spf13/cobra"
)

var (
	ListCmd = &cobra.Command{
		Use:   "list",
		Short: "List crusado profiles",
		Long:  `Allows you to list available profiles. You can specify an output format.`,
		Args:  cobra.NoArgs,
		Run:   List,
	}
)

func init() {
	ListCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "", "define the output format: [wide, yaml, json, jsonpath]")
}

func List(cmd *cobra.Command, args []string) {
	// TODO read the list of availble profiles from a well-known source, e.g. $HOME/.crusado/profiles.yaml
	// currently we only support one profile
	cfg := config.GetConfig()

	profileConfigs := []config.ProfileConfig{
		{
			Name:     "default",
			FilePath: cfg.ProfileFilePath,
		},
	}

	profiles := []config.Profile{}

	for i := range profileConfigs {
		profiles = append(profiles, convertToProfile(&profileConfigs[i]))
	}

	printer, err := getPrinter(outputFlag)
	if err != nil {
		log.Fatalf("Could not get printer: %q", err)
	}

	printer.Fprint(os.Stdout, profiles)
	if err != nil {
		log.Fatalf("Could not print profiles: %q", err)
	}
}

func convertToProfile(profileConfig *config.ProfileConfig) config.Profile {
	profile := &config.Profile{}

	profile.Name = profileConfig.Name
	profile.FilePath = profileConfig.FilePath

	errs := validator.ValidateProfileConfig(profileConfig)

	if len(errs) != 0 {
		profile.Valid = false
		return *profile
	}

	templateList, err := config.GetTemplateListFromFile(profile.FilePath)
	if err != nil {
		profile.Valid = false
		return *profile
	}

	err = validator.ValidateTemplateList(templateList)
	if err != nil {
		profile.Valid = false
		return *profile
	}

	profile.Valid = true
	profile.NumberOfTemplates = len(templateList.Templates)

	return *profile
}
