package templates

import (
	"errors"
	"fmt"
)

var (
	ErrDuplicateTemplateNames = errors.New("duplicate names for templates")
	ErrTypeNotSet             = errors.New("template doesn't have type, which is required")
	ErrInvalidType            = errors.New("specified type is not valid")
)

// ValidateTemplateList validates the list of templates given as a whole as well
// as each indiviual template within the list. It returns an error that is
// constructed using errors.Join().
func ValidateTemplateList(templateList *TemplateList) error {
	errs := []error{}

	errs = append(errs, ValidateTemplateListUniqueName(templateList))

	for i := range templateList.Templates {
		errs = append(errs, ValidateTemplate(&templateList.Templates[i]))
	}

	return errors.Join(errs...)
}

func ValidateTemplateListUniqueName(templateList *TemplateList) error {
	templateNames := map[string]bool{}
	errs := []error{}

	for i := range templateList.Templates {
		name := templateList.Templates[i].Name
		if exists := templateNames[name]; exists {
			errs = append(errs, fmt.Errorf("%w: name '%s' exists at least twice", ErrDuplicateTemplateNames, name))
		} else {
			templateNames[name] = true
		}
	}

	return errors.Join(errs...)
}

// ValidateTemplate validates the given template. It returns an error that is
// constructed using errors.Join().
func ValidateTemplate(template *Template) error {
	var errs []error

	errs = append(errs, ValidateTemplateValidType(template))

	return errors.Join(errs...)
}

func ValidateTemplateValidType(template *Template) error {
	if template.Type == "" {
		return ErrTypeNotSet
	}

	for i := range AvailableTypes {
		if template.Type == AvailableTypes[i] {
			return nil
		}
	}

	return fmt.Errorf("%w: type '%s' should be one of %+v", ErrInvalidType, template.Type, AvailableTypes)
}
