package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewConsultationCase(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	title, _ := NewCaseTitle("Contract Review")
	description, _ := NewCaseDescription("Need advice reviewing a contract")
	category, _ := NewCaseCategory("LEGAL")

	t.Run("should successfully create a new consultation case", func(t *testing.T) {
		c, err := NewConsultationCase(
			"case_123",
			"user_456",
			*title,
			*description,
			*category,
			now,
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		require.Equal(t, "case_123", c.ID())
		require.Equal(t, "user_456", c.ClientID())
		require.Equal(t, "Contract Review", c.Title().String())
		require.Equal(t, "Need advice reviewing a contract", c.Description().String())
		require.Equal(t, "LEGAL", c.Category().String())
		require.Equal(t, CaseStatusDraft, c.Status())
		require.Equal(t, now, c.CreatedAt())
		require.Equal(t, now, c.UpdatedAt())
	})

	t.Run("should fail when ID is empty", func(t *testing.T) {
		c, err := NewConsultationCase(
			"",
			"user_456",
			*title,
			*description,
			*category,
			now,
		)
		require.Error(t, err)
		require.Nil(t, c)
		require.Equal(t, ErrInvalidCaseID, err)
	})

	t.Run("should fail when ClientID is empty", func(t *testing.T) {
		c, err := NewConsultationCase(
			"case_123",
			"",
			*title,
			*description,
			*category,
			now,
		)
		require.Error(t, err)
		require.Nil(t, c)
		require.Equal(t, ErrInvalidClientID, err)
	})
}

func TestReconstituteConsultationCase(t *testing.T) {
	createdAt := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC)

	t.Run("should successfully reconstitute consultation case", func(t *testing.T) {
		c, err := ReconstituteConsultationCase(
			"case_123",
			"user_456",
			"Title",
			"Description",
			"TECH",
			"SUBMITTED",
			createdAt,
			updatedAt,
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		require.Equal(t, "case_123", c.ID())
		require.Equal(t, "user_456", c.ClientID())
		require.Equal(t, "Title", c.Title().String())
		require.Equal(t, "Description", c.Description().String())
		require.Equal(t, "TECH", c.Category().String())
		require.Equal(t, CaseStatusSubmitted, c.Status())
		require.Equal(t, createdAt, c.CreatedAt())
		require.Equal(t, updatedAt, c.UpdatedAt())
	})

	t.Run("should fail with empty case ID", func(t *testing.T) {
		_, err := ReconstituteConsultationCase("", "user_456", "Title", "Desc", "TECH", "DRAFT", createdAt, updatedAt)
		require.Equal(t, ErrInvalidCaseID, err)
	})

	t.Run("should fail with empty client ID", func(t *testing.T) {
		_, err := ReconstituteConsultationCase("case_123", "", "Title", "Desc", "TECH", "DRAFT", createdAt, updatedAt)
		require.Equal(t, ErrInvalidClientID, err)
	})

	t.Run("should fail with invalid title", func(t *testing.T) {
		_, err := ReconstituteConsultationCase("case_123", "user_456", "", "Desc", "TECH", "DRAFT", createdAt, updatedAt)
		require.Equal(t, ErrCaseTitleEmpty, err)
	})

	t.Run("should fail with invalid description", func(t *testing.T) {
		_, err := ReconstituteConsultationCase("case_123", "user_456", "Title", "", "TECH", "DRAFT", createdAt, updatedAt)
		require.Equal(t, ErrCaseDescriptionEmpty, err)
	})

	t.Run("should fail with invalid category", func(t *testing.T) {
		_, err := ReconstituteConsultationCase("case_123", "user_456", "Title", "Desc", "", "DRAFT", createdAt, updatedAt)
		require.Equal(t, ErrCaseCategoryEmpty, err)
	})

	t.Run("should fail with invalid status", func(t *testing.T) {
		_, err := ReconstituteConsultationCase("case_123", "user_456", "Title", "Desc", "TECH", "INVALID", createdAt, updatedAt)
		require.Equal(t, ErrCaseStatusInvalid, err)
	})
}

func TestConsultationCaseUpdatesAndTransitions(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	title, _ := NewCaseTitle("Initial Title")
	description, _ := NewCaseDescription("Initial Description")
	category, _ := NewCaseCategory("LEGAL")

	c, err := NewConsultationCase("case_123", "user_456", *title, *description, *category, now)
	require.NoError(t, err)

	t.Run("update title, description, category, and end updating", func(t *testing.T) {
		newTitle, _ := NewCaseTitle("Updated Title")
		newDesc, _ := NewCaseDescription("Updated Description")
		newCat, _ := NewCaseCategory("TAX")
		updateTime := now.Add(time.Hour)

		c.UpdateTitle(*newTitle)
		c.UpdateDescription(*newDesc)
		c.UpdateCategory(*newCat)
		c.EndUpdating(updateTime)

		require.Equal(t, "Updated Title", c.Title().String())
		require.Equal(t, "Updated Description", c.Description().String())
		require.Equal(t, "TAX", c.Category().String())
		require.Equal(t, updateTime, c.UpdatedAt())
	})

	t.Run("valid status transition", func(t *testing.T) {
		transitionTime := now.Add(2 * time.Hour)
		err := c.TransitionToStatus(CaseStatusSubmitted, transitionTime)
		require.NoError(t, err)
		require.Equal(t, CaseStatusSubmitted, c.Status())
		require.Equal(t, transitionTime, c.UpdatedAt())
	})

	t.Run("invalid status transition", func(t *testing.T) {
		invalidTime := now.Add(3 * time.Hour)
		err := c.TransitionToStatus(CaseStatusDraft, invalidTime)
		require.Error(t, err)
		require.Equal(t, ErrInvalidStatusTransition, err)
	})
}
