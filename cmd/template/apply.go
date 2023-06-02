package template

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/simonkienzler/crusado/pkg/templates"
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
		Short: "Apply a user story or bug with tasks based on crusado templates",
		Long: `Allows you to create a user story or bug from the template specified by the argument
given to the command. Able to execute in dry-run mode, if you don't actually
want to create any workitems in Azure DevOps.`,
		Args: cobra.ExactArgs(1),
		Run:  Apply,
	}
)

var (
	dryRunFlag          bool
	iterationOffsetFlag int
)

func init() {
	// TODO change the default to false at some point, it's the more sensible default for actual usage
	ApplyCmd.PersistentFlags().BoolVarP(&dryRunFlag, "dry-run", "d", true, "if set to true, crusado doesn't actually create work items in Azure DevOps")
	ApplyCmd.PersistentFlags().IntVarP(&iterationOffsetFlag, "iteration-offset", "i", 1, "iteration to apply the template in, relative to the current iteration.\n1 will traget the next iteration, -1 the previous one.")
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
		prettyPrintValidationError(err)
	}

	ustService := &templates.Service{
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

	iteration, err := workitemsService.GetIterationRelativeToCurrent(ctx, iterationOffsetFlag)
	if err != nil {
		return nil, err
	}

	workitemsService.ProjectConfig.IterationPath = *iteration.Path

	return &workitemsService, nil
}

func ApplyFlow(ctx context.Context, service *templates.Service, templateName string, dryRun bool) {
	template, err := service.GetTemplateFromName(templateName)
	if err != nil {
		log.Fatalf("Could not get template: %s", err)
	}

	storyDescriptionHTML, err := templates.ConvertMarkdownToHTML(template.Description)
	if err != nil {
		log.Fatalf("Could not convert story desription from Markdown to HTML: %s", err)
	}

	userStory, err := service.WorkitemsService.Create(ctx, template.Title, storyDescriptionHTML, template.Type)
	if err != nil {
		log.Fatalf("Could not create from template '%s': %s", templateName, err)
	}

	coloredIterationPathPrinter(service.WorkitemsService.ProjectConfig.IterationPath)
	coloredSuccessMessagePrinter(template.Type, template.Title, dryRun)

	for i := range template.Tasks {
		task := template.Tasks[i]
		taskDescriptionHTML, err := templates.ConvertMarkdownToHTML(task.Description)
		if err != nil {
			log.Fatalf("Could not convert task desription from Markdown to HTML: %s", err)
		}

		_, err = service.WorkitemsService.CreateTaskUnderneath(ctx, task.Title, taskDescriptionHTML, userStory)
		if err != nil {
			log.Fatalf("Could not create task: %s", err)
		}

		coloredSuccessMessagePrinter(workitems.TaskType, task.Title, dryRun)
	}
}

func coloredSuccessMessagePrinter(templateType config.TemplateType, title string, dryRun bool) {
	const (
		storyIcon = "📖"
		bugIcon   = "🐛"
		taskIcon  = "📋"
	)

	dryRunHint := ""

	if dryRun {
		dryRunHint = "(dry-run)"
	}

	icon := ""
	itemType := ""
	txtColor := color.FgCyan

	switch templateType {
	case config.TemplateTypeUserStory:
		icon = storyIcon
		txtColor = color.FgGreen
		itemType = workitems.UserStoryType
	case config.TemplateTypeBug:
		icon = bugIcon
		txtColor = color.FgRed
		itemType = workitems.BugType
	case workitems.TaskType:
		icon = "   " + taskIcon
		itemType = workitems.TaskType
	}

	fmt.Print(icon + " " + itemType + " ")
	color.New(txtColor).Print(title)
	fmt.Printf(" created successfully %s\n", dryRunHint)
}

func coloredIterationPathPrinter(iterationPath string) {
	const (
		iterationIcon = "🔁"
	)

	parts := strings.Split(iterationPath, "\\")

	fmt.Print(iterationIcon + " Iteration Path: ")

	for i := range parts {
		color.New(color.FgYellow).Print(parts[i])

		if i < len(parts)-1 {
			fmt.Print(" > ")
		}
	}

	fmt.Print("\n")
}
