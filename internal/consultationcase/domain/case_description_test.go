package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCaseDescription(t *testing.T) {
	t.Run("should create valid case description", func(t *testing.T) {
		desc, err := NewCaseDescription("Detailed description of the issue encountered.")
		require.NoError(t, err)
		require.NotNil(t, desc)
		require.Equal(t, "Detailed description of the issue encountered.", desc.String())
	})

	t.Run("should succeed for 500 character description", func(t *testing.T) {
		longDesc := strings.Repeat("b", 500)
		desc, err := NewCaseDescription(longDesc)
		require.NoError(t, err)
		require.Equal(t, longDesc, desc.String())
	})

	t.Run("should fail when description is empty", func(t *testing.T) {
		desc, err := NewCaseDescription("")
		require.Error(t, err)
		require.Nil(t, desc)
		require.Equal(t, ErrCaseDescriptionEmpty, err)
	})

	t.Run("should fail when description exceeds 500 characters", func(t *testing.T) {
		tooLongDesc := strings.Repeat("b", 501)
		desc, err := NewCaseDescription(tooLongDesc)
		require.Error(t, err)
		require.Nil(t, desc)
		require.Equal(t, ErrCaseDescriptionTooLong, err)
	})
}
