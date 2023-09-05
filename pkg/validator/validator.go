package validator

import (
	"errors"
	"fmt"
	"os"

	"github.com/simonkienzler/crusado/pkg/config"
)

var (
	ErrProfileNameNotSet      = errors.New("profile doesn't have name, which is required")
	ErrProfileFilePathNotSet  = errors.New("profile doesn't have filePath, which is required")
	ErrFileDoesNotExist       = errors.New("no file exists at the given filePath")
	ErrDuplicateTemplateNames = errors.New("duplicate names for templates in profile")
	ErrTypeNotSet             = errors.New("template doesn't have type, which is required")
	ErrInvalidType            = errors.New("specified type is not valid")
)

func ValidateProfileConfig(profile *config.ProfileConfig) []error {
	validationErrors := []error{}

	// validate completeness
	if profile.Name == "" {
		validationErrors = append(validationErrors, ErrProfileNameNotSet)
	}

	if profile.FilePath == "" {
		validationErrors = append(validationErrors, ErrProfileFilePathNotSet)
		// we can return immediately because there is no file to check
		return validationErrors
	}

	// validate existence and conformance of profile file
	if _, err := os.Stat(profile.FilePath); errors.Is(err, os.ErrNotExist) {
		validationErrors = append(validationErrors, ErrFileDoesNotExist)
		return validationErrors
	}

	return validationErrors
}

// ValidateTemplateList validates the list of templates given as a whole as well
// as each indiviual template within the list. It returns an error that is
// constructed using errors.Join().
func ValidateTemplateList(templateList *config.TemplateList) error {
	errs := []error{}

	errs = append(errs, ValidateTemplateListUniqueName(templateList))

	for i := range templateList.Templates {
		errs = append(errs, ValidateTemplate(&templateList.Templates[i]))
	}

	return errors.Join(errs...)
}

func ValidateTemplateListUniqueName(templateList *config.TemplateList) error {
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
func ValidateTemplate(template *config.Template) error {
	var errs []error

	errs = append(errs, ValidateTemplateValidType(template))

	return errors.Join(errs...)
}

func ValidateTemplateValidType(template *config.Template) error {
	if template.Type == "" {
		return ErrTypeNotSet
	}

	for i := range config.AvailableTemplateTypes {
		if template.Type == config.AvailableTemplateTypes[i] {
			return nil
		}
	}

	return fmt.Errorf("%w: type '%s' should be one of %+v", ErrInvalidType, template.Type, config.AvailableTemplateTypes)
}
