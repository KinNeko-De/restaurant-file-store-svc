package file

import (
	context "context"
	"testing"

	"github.com/google/uuid"
	ioFixture "github.com/kinneko-de/restaurant-file-store-svc/test/testing/io"
	mock "github.com/stretchr/testify/mock"
)

func (mockFileRepository *MockFileRepository) setupCreateFileNewFile(t *testing.T, fileWriter *ioFixture.MockWriteCloser) func() (uuid.UUID, uuid.UUID) {
	t.Helper()
	var generatedFileId uuid.UUID
	var generatedRevisionId uuid.UUID
	mockFileRepository.EXPECT().CreateFile(mock.Anything, mock.IsType(uuid.New()), mock.IsType(uuid.New())).
		Run(func(ctx context.Context, fileId uuid.UUID, revisionId uuid.UUID) {
			generatedFileId = fileId
			generatedRevisionId = revisionId
		}).
		Return(fileWriter, nil).
		Times(1)

	return func() (uuid.UUID, uuid.UUID) {
		return generatedFileId, generatedRevisionId
	}
}

func (mockFileRepository *MockFileRepository) setupCreateFileError(t *testing.T, createError error) {
	t.Helper()
	mockFileRepository.EXPECT().CreateFile(mock.Anything, mock.IsType(uuid.New()), mock.IsType(uuid.New())).Return(nil, createError).Times(1)
}

func (mockFileRepository *MockFileRepository) setupCreateFileNewRevision(t *testing.T, expectedFileId uuid.UUID, fileWriter *ioFixture.MockWriteCloser) func() uuid.UUID {
	t.Helper()
	var generatedRevisionId uuid.UUID
	mockFileRepository.EXPECT().CreateFile(mock.Anything, expectedFileId, mock.IsType(uuid.New())).
		Run(func(ctx context.Context, fileId uuid.UUID, revisionId uuid.UUID) {
			generatedRevisionId = revisionId
		}).
		Return(fileWriter, nil).
		Times(1)

	return func() uuid.UUID {
		return generatedRevisionId
	}
}
