package crusado

import "github.com/thediveo/klo"

type Template struct {
	Meta        Meta
	Description string
}

// Meta represents the Metadata associated with a Crusado template
type Meta struct {
	// Name is the unique name of the template, used in commands
	Name string `yaml:"name" json:"name"`

	// Summary provides a short synopsis for the template
	Summary string `yaml:"summary" json:"summary"`

	// Type identifies the template as one of [User Story, Bug]
	Type Type `yaml:"type" json:"type"`

	// Title is the resulting title of the work item in Azure DevOps
	Title string `yaml:"title" json:"title"`

	// Tasks is a slice of individual tasks that are part of the template
	Tasks []Task `yaml:"tasks" json:"tasks"`
}

type Task struct {
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description" json:"description"`
}

type Type string

const (
	UserStoryType = Type("UserStory")
	BugType       = Type("Bug")
	TaskType      = Type("Task")
)

var AvailableTypes = []Type{
	UserStoryType,
	BugType,
}

var PrinterSpecs = klo.Specs{
	DefaultColumnSpec: "NAME:{.Meta.Name},TYPE:{.Meta.Type},SUMMARY:{.Meta.Summary}",
	WideColumnSpec:    "NAME:{.Meta.Name},TYPE:{.Meta.Type},SUMMARY:{.Meta.Summary},TITLE:{.Meta.Title},TASKS:{.Meta.Tasks[*].Title}",
}
