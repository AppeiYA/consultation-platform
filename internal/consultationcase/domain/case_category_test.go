package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCaseCategory(t *testing.T) {
	t.Run("should create valid case category", func(t *testing.T) {
		cat, err := NewCaseCategory("FINANCE")
		require.NoError(t, err)
		require.NotNil(t, cat)
		require.Equal(t, "FINANCE", cat.String())
	})

	t.Run("should fail when category is empty", func(t *testing.T) {
		cat, err := NewCaseCategory("")
		require.Error(t, err)
		require.Nil(t, cat)
		require.Equal(t, ErrCaseCategoryEmpty, err)
	})
}
