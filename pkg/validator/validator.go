package validator

import (
	"errors"
	"os"

	"github.com/simonkienzler/crusado/pkg/config"
)

var (
	ErrProfileNameNotSet     = errors.New("profile doesn't have name, which is required")
	ErrProfileFilePathNotSet = errors.New("profile doesn't have filePath, which is required")
	ErrFileDoesNotExist      = errors.New("no file exists at the given filePath")
)

func ValidateProfileConfig(profile config.ProfileConfig) []error {
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
