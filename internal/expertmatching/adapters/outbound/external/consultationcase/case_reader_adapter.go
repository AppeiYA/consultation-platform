package consultationcase_adapter

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/shared/db"
	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

type CaseReaderAdapter struct {
	repository *db.Repository
}

func NewCaseReaderAdapter(repository *db.Repository) *CaseReaderAdapter {
	return &CaseReaderAdapter{
		repository: repository,
	}
}

type caseRow struct {
	ID          string `db:"id"`
	ClientID    string `db:"client_id"`
	Category    string `db:"category"`
	Title       string `db:"title"`
	Description string `db:"description"`
}

func (a *CaseReaderAdapter) GetCaseDetails(ctx context.Context, caseID string) (*outbound.CaseDetails, error) {
	query := `
		SELECT id, client_id, category, title, description
		FROM consultation_cases
		WHERE id = $1
	`
	var row caseRow
	err := a.repository.Executor(ctx).GetContext(ctx, &row, query, caseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, custom_errors.NotFoundError("consultation case not found")
		}
		return nil, err
	}

	category, err := domain.NewMatchingCategory(row.Category)
	if err != nil {
		return nil, err
	}

	return &outbound.CaseDetails{
		CaseID:      row.ID,
		ClientID:    row.ClientID,
		Category:    category,
		Title:       row.Title,
		Description: row.Description,
	}, nil
}

var _ outbound.CaseReader = (*CaseReaderAdapter)(nil)
