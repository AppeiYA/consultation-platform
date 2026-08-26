package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCaseTitle(t *testing.T) {
	t.Run("should create valid case title", func(t *testing.T) {
		title, err := NewCaseTitle("Legal Advice Needed")
		require.NoError(t, err)
		require.NotNil(t, title)
		require.Equal(t, "Legal Advice Needed", title.String())
	})

	t.Run("should succeed for 100 character title", func(t *testing.T) {
		longTitle := strings.Repeat("a", 100)
		title, err := NewCaseTitle(longTitle)
		require.NoError(t, err)
		require.Equal(t, longTitle, title.String())
	})

	t.Run("should fail when title is empty", func(t *testing.T) {
		title, err := NewCaseTitle("")
		require.Error(t, err)
		require.Nil(t, title)
		require.Equal(t, ErrCaseTitleEmpty, err)
	})

	t.Run("should fail when title exceeds 100 characters", func(t *testing.T) {
		tooLongTitle := strings.Repeat("a", 101)
		title, err := NewCaseTitle(tooLongTitle)
		require.Error(t, err)
		require.Nil(t, title)
		require.Equal(t, ErrCaseTitleTooLong, err)
	})
}
