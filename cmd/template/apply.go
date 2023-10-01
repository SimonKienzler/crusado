package template

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/simonkienzler/crusado/pkg/crusado"
	"github.com/simonkienzler/crusado/pkg/workitems"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	ApplyCmd = &cobra.Command{
		Use:   "apply [template-name]",
		Short: "Apply a user story or bug with tasks based on crusado templates",
		Long: `Allows you to create a user story or bug from the template specified by the argument
given to the command. Supports dry-run, skipping confirmation and let's you specify
the exact iteration in which to apply the template.`,
		Args: cobra.ExactArgs(1),
		Run:  Apply,
	}
)

var (
	dryRunFlag          bool
	iterationOffsetFlag int
	autoApproveFlag     bool
)

func init() {
	ApplyCmd.PersistentFlags().BoolVarP(&autoApproveFlag, "yes", "y", false, "skip confirmation step")

	dryRunDesc := "if set to true, crusado doesn't actually create work items in Azure DevOps"
	ApplyCmd.PersistentFlags().BoolVarP(&dryRunFlag, "dry-run", "d", false, dryRunDesc)

	iterationOffsetDesc := "iteration to apply the template in, relative to the current iteration.\n1 will traget the next iteration, -1 the previous one."
	ApplyCmd.PersistentFlags().IntVarP(&iterationOffsetFlag, "iteration-offset", "i", 1, iterationOffsetDesc)
}

func Apply(_ *cobra.Command, args []string) {
	// TODO implement proper contexts
	ctx := context.Background()

	wiService, err := workitemsService(ctx, dryRunFlag)
	if err != nil {
		log.Fatalf("Error during service creation: %s", err)
	}

	ApplyFlow(ctx, crusadoService(), wiService, args[0])
}

func ApplyFlow(ctx context.Context, tplService *crusado.Service, wiService *workitems.Service, templateName string) {
	var createdItemHint string

	if dryRunFlag {
		createdItemHint = "would be created"
	} else {
		createdItemHint = "created successfully"
	}

	template, err := tplService.GetByName(templateName)
	if err != nil {
		log.Fatalf("Could not get template: %s", err)
	}

	coloredIterationPathPrinter(wiService.IterationPath)

	if !autoApproveFlag {
		coloredItemPrinter(template.Meta.Type, template.Meta.Title, "")

		for i := range template.Meta.Tasks {
			task := template.Meta.Tasks[i]
			coloredItemPrinter(workitems.TaskType, task.Title, "")
		}

		if !confirm("Create these work items in the specified iteration path?") {
			fmt.Printf("No work items created.\n")
			return
		}

		fmt.Println()
	}

	userStory, err := wiService.Create(ctx, template.Meta.Title, template.Description, template.Meta.Type)
	if err != nil {
		log.Fatalf("Could not create from template '%s': %s", templateName, err)
	}

	createdItemHintWithURL := createdItemHint

	url, err := wiService.GetWorkItemHTMLRef(userStory)
	if err == nil && url != nil {
		createdItemHintWithURL += fmt.Sprintf(" at %s", *url)
	}

	coloredItemPrinter(template.Meta.Type, template.Meta.Title, createdItemHintWithURL)

	for i := range template.Meta.Tasks {
		task := template.Meta.Tasks[i]

		_, err = wiService.CreateTaskUnderneath(ctx, task.Title, task.Description, userStory)
		if err != nil {
			log.Fatalf("Could not create task: %s", err)
		}

		coloredItemPrinter(workitems.TaskType, task.Title, createdItemHint)
	}
}

func coloredItemPrinter(templateType crusado.Type, title, addendum string) {
	const (
		storyIcon = "📖"
		bugIcon   = "🐛"
		taskIcon  = "📋"
	)

	icon := ""
	itemType := ""
	txtColor := color.FgCyan

	switch templateType {
	case crusado.UserStoryType:
		icon = storyIcon
		txtColor = color.FgGreen
		itemType = workitems.UserStoryType
	case crusado.BugType:
		icon = bugIcon
		txtColor = color.FgRed
		itemType = workitems.BugType
	case workitems.TaskType:
		icon = "   " + taskIcon
		itemType = workitems.TaskType
	}

	fmt.Print(icon + " " + itemType + " ")
	color.New(txtColor).Print(title)
	fmt.Printf(" %s\n", addendum)
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

	fmt.Print("\n\n")
}

func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("\n%s [y/n]: ", prompt)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
