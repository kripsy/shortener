package utils

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveCert(t *testing.T) {
	t.Run("save valid cert", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "example.*.txt")
		assert.NoError(t, err)
		defer os.Remove(tmpfile.Name())

		err = tmpfile.Close()
		assert.NoError(t, err)

		expectedContent := "test content"
		payload := bytes.NewBufferString(expectedContent)

		err = saveCert(tmpfile.Name(), payload)
		assert.NoError(t, err)

		content, err := os.ReadFile(tmpfile.Name())
		assert.NoError(t, err)
		assert.Equal(t, expectedContent, string(content))
	})
}
