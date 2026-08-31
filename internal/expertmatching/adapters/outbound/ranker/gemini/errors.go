package geminiranker

import (
	"fmt"

	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

var (
	ErrGeneratingCandidateRanking = func(message string, err error) error {
		genErr := fmt.Sprintf("%s: %v", message, err)

		return custom_errors.InternalServerError(genErr)
	}
)