package template

import (
	"context"
	"fmt"
	"log"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/simonkienzler/crusado/pkg/userstorytemplates"
	"github.com/simonkienzler/crusado/pkg/validator"
	"github.com/simonkienzler/crusado/pkg/workitems"

	"github.com/fatih/color"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/work"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/spf13/cobra"
)

var (
	ApplyCmd = &cobra.Command{
		Use:   "apply [template-name]",
		Short: "Apply a user story with tasks based on crusado templates",
		Long: `Allows you to create a user story from the template specified by the argument
given to the command. Able to execute in dry-run mode, if you don't actually
want to create any workitems in Azure DevOps.`,
		Args: cobra.ExactArgs(1),
		Run:  Apply,
	}
)

var (
	dryRunFlag bool
)

func init() {
	// TODO change the default to false at some point, it's the more sensible default for actual usage
	ApplyCmd.PersistentFlags().BoolVar(&dryRunFlag, "dry-run", true, "if set to true, crusado doesn't actually create work items in Azure DevOps")
}

func Apply(cmd *cobra.Command, args []string) {
	// TODO implement proper contexts
	ctx := context.Background()

	cfg := config.GetConfig()

	workitemsService, err := createWorkitemsService(ctx, dryRunFlag)
	if err != nil {
		log.Fatalf("Error during service creation: %s", err)
	}

	// create templateList from example
	templateList, err := config.GetTemplateListFromFile(cfg.ProfileFilePath)
	if err != nil {
		log.Fatalf("Could not read example template file: %s", err)
	}

	err = validator.ValidateTemplateList(templateList)
	if err != nil {
		log.Fatalf("Invalid profile: %s", err)
	}

	ustService := &userstorytemplates.Service{
		WorkitemsService: *workitemsService,
		TemplateList:     *templateList,
	}

	ApplyFlow(ctx, ustService, args[0], dryRunFlag)
}

func createWorkitemsService(ctx context.Context, useDryRunMode bool) (*workitems.Service, error) {
	crusadoConfig := config.GetConfig()

	// create a connection to the organization
	connection := azuredevops.NewPatConnection(crusadoConfig.OrganizationUrl, crusadoConfig.PersonalAccessToken)

	// create clients required by the workitems service
	workitemClient, err := workitemtracking.NewClient(ctx, connection)
	if err != nil {
		return nil, err
	}

	workClient, err := work.NewClient(ctx, connection)
	if err != nil {
		return nil, err
	}

	// configure the workitems service
	workitemsService := workitems.Service{
		WorkitemClient: workitemClient,
		WorkClient:     workClient,
		DryRun:         useDryRunMode,
		ProjectConfig: config.ProjectConfig{
			Name:     crusadoConfig.ProjectName,
			AreaPath: crusadoConfig.ProjectName,
		},
	}

	// get current iteration, either from env var or from the API
	workitemsService.ProjectConfig.IterationPath = crusadoConfig.IterationPath
	if crusadoConfig.UseCurrentIteration {
		currentIteration, err := workitemsService.GetCurrentIteration(ctx)
		if err != nil {
			return nil, err
		}

		log.Printf("Got current iteration path: %+v", *currentIteration.Path)
		workitemsService.ProjectConfig.IterationPath = *currentIteration.Path
	}

	return &workitemsService, nil
}

func ApplyFlow(ctx context.Context, service *userstorytemplates.Service, templateName string, dryRun bool) {
	userStoryTemplate, err := service.GetUserStoryTemplateFromName(templateName)
	if err != nil {
		log.Fatalf("Could not get user story template: %s", err)
	}

	storyDescriptionHTML, err := userstorytemplates.ConvertMarkdownToHTML(userStoryTemplate.StoryDescription)
	if err != nil {
		log.Fatalf("Could not convert story desription from Markdown to HTML: %s", err)
	}

	userStory, err := service.WorkitemsService.CreateUserStory(ctx, userStoryTemplate.StoryTitle, storyDescriptionHTML)
	if err != nil {
		log.Fatalf("Could not create user story: %s", err)
	}

	coloredSuccessMessagePrinter(workitems.UserStoryType, userStoryTemplate.StoryTitle, dryRun)

	for i := range userStoryTemplate.Tasks {
		task := userStoryTemplate.Tasks[i]
		taskDescriptionHTML, err := userstorytemplates.ConvertMarkdownToHTML(userStoryTemplate.StoryDescription)
		if err != nil {
			log.Fatalf("Could not convert task desription from Markdown to HTML: %s", err)
		}

		_, err = service.WorkitemsService.CreateTaskUnderneathUserStory(ctx, task.Title, taskDescriptionHTML, userStory)
		if err != nil {
			log.Fatalf("Could not create task: %s", err)
		}

		coloredSuccessMessagePrinter(workitems.TaskType, task.Title, dryRun)
	}
}

func coloredSuccessMessagePrinter(itemType, title string, dryRun bool) {
	const (
		storyIcon = "📖"
		bugIcon   = "🐛"
		taskIcon  = "📋"
	)

	dryRunHint := ""

	if dryRun {
		dryRunHint = " (dry-run)"
	}

	icon := ""
	txtColor := color.FgCyan

	switch itemType {
	case workitems.UserStoryType:
		icon = storyIcon
		txtColor = color.FgGreen
	case workitems.TaskType:
		icon = "   " + taskIcon
	}

	fmt.Print(icon + " " + itemType + " ")
	color.New(txtColor).Printf(title)
	fmt.Printf(" created successfully%s\n", dryRunHint)
}
