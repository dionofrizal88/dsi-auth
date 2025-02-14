package util_test

import (
	"github.com/dionofrizal88/dsi/auth/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestUtilExtractor(t *testing.T) {
	byteFileTest, _ := os.ReadFile("../pkg/tests/file-test-1.png")

	t.Run("positive case to test extractor utility, expected no error", func(t *testing.T) {
		t.Run("positive case while use func get type utility, expected no error", func(t *testing.T) {

			byteData, fileType, err := util.GetType(byteFileTest)
			require.NoError(t, err)
			assert.NotEmpty(t, len(byteData))
			assert.Equal(t, "image/png", fileType)
		})
	})

}
