//go:build unit

package persistence

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fixture "github.com/kinneko-de/restaurant-file-store-svc/test/testing/file"
)

func TestCreateFile(t *testing.T) {
	fileId := uuid.New()
	revisionId := uuid.New()
	rootPath := t.TempDir()
	expectedFileAndPath := path.Join(rootPath, fileId.String(), revisionId.String())
	expectedContent := fixture.TextFile()

	sut := PersistentVolumeFileRepository{StoragePath: rootPath}

	writer, err := sut.CreateFile(context.Background(), fileId, revisionId, 1024)
	require.Nil(t, err)
	assert.NotNil(t, writer)
	writer.Write(expectedContent)
	writer.Close()

	content, err := os.ReadFile(expectedFileAndPath)
	require.Nil(t, err)
	assert.Equal(t, expectedContent, content)
}
