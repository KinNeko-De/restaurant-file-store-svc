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
	writer, err := sut.CreateFile(context.Background(), fileId, revisionId)

	require.Nil(t, err)
	assert.NotNil(t, writer)
	writer.Write(expectedContent)
	writer.Close()

	content, err := os.ReadFile(expectedFileAndPath)
	require.Nil(t, err)
	assert.Equal(t, expectedContent, content)
}

func TestOpenFile(t *testing.T) {
	fileId := uuid.New()
	revisionId := uuid.New()
	rootPath := t.TempDir()
	expectedContent := fixture.TextFile()
	buffer := make([]byte, len(expectedContent)+10)

	fileAndPath := path.Join(rootPath, fileId.String(), revisionId.String())
	_ = os.MkdirAll(path.Dir(fileAndPath), os.ModePerm)
	writer, err := os.Create(fileAndPath)
	writer.Write(expectedContent)
	writer.Close()

	sut := PersistentVolumeFileRepository{StoragePath: rootPath}
	reader, err := sut.OpenFile(context.Background(), fileId, revisionId)
	defer reader.Close()

	require.Nil(t, err)
	assert.NotNil(t, reader)
	readBytes, err := reader.Read(buffer)
	require.Nil(t, err)
	assert.Equal(t, len(expectedContent), readBytes)
	assert.Equal(t, expectedContent, buffer[:readBytes])
}
