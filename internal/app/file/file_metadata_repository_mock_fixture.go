package file

import (
	context "context"
	"errors"
	"testing"

	"github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

func (mockFileMetadataRepository *MockFileMetadataRepository) setupStoreFileMetadata(t *testing.T) func() FileMetadata {
	t.Helper()
	var storedFileMetadata FileMetadata
	mockFileMetadataRepository.EXPECT().StoreFileMetadata(mock.Anything, mock.IsType(FileMetadata{})).
		Run(func(ctx context.Context, fileMetadata FileMetadata) { storedFileMetadata = fileMetadata }).
		Return(nil).
		Times(1)

	return func() FileMetadata {
		return storedFileMetadata
	}
}

func (mockFileMetadataRepository *MockFileMetadataRepository) setupStoreRevisionMetadata(t *testing.T, fileid uuid.UUID) func() Revision {
	t.Helper()
	var storedRevisionMetadata Revision
	mockFileMetadataRepository.EXPECT().StoreRevision(mock.Anything, fileid, mock.IsType(Revision{})).
		Run(func(ctx context.Context, existingFileId uuid.UUID, revision Revision) {
			storedRevisionMetadata = revision
		}).
		Return(nil).
		Times(1)

	return func() Revision {
		return storedRevisionMetadata
	}
}

func (mockFileMetadataRepository *MockFileMetadataRepository) setupFileMetadataRepositoryMockStoreFileMetadataReturnsError(t *testing.T, err error) *MockFileMetadataRepository {
	t.Helper()
	mockFileMetadataRepository.EXPECT().StoreFileMetadata(mock.Anything, mock.IsType(FileMetadata{})).
		Return(err).
		Times(1)

	return mockFileMetadataRepository
}

func (mockFileMetadataRepository *MockFileMetadataRepository) setupFileMetadataRepositoryMockStoreRevisionReturnsError(t *testing.T, fileId uuid.UUID, err error) {
	t.Helper()
	mockFileMetadataRepository.EXPECT().StoreRevision(mock.Anything, fileId, mock.IsType(Revision{})).
		Return(err).
		Times(1)
}

func (mockFileMetadataRepository *MockFileMetadataRepository) setupFileMetadataRepositoryToFetchMetadata(t *testing.T, fileId uuid.UUID, fileMetadata FileMetadata) {
	t.Helper()
	mockFileMetadataRepository.EXPECT().FetchFileMetadata(mock.Anything, fileId).Return(fileMetadata, nil).Times(1)
	mockFileMetadataRepository.EXPECT().NotFoundError().Return(errors.New("not expected error")).Times(1)
}
