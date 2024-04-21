package file

import (
	context "context"
	"errors"
	"testing"

	"github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

func setupStoreFileMetadata(t *testing.T, mockFileMetadataRepository *MockFileMetadataRepository, storedFileMetadata **FileMetadata) {
	t.Helper()
	mockFileMetadataRepository.EXPECT().StoreFileMetadata(mock.Anything, mock.IsType(FileMetadata{})).
		Run(func(ctx context.Context, fileMetadata FileMetadata) { *storedFileMetadata = &fileMetadata }).
		Return(nil).
		Times(1)
}

func setupStoreRevisionMetadata(t *testing.T, mockFileMetadataRepository *MockFileMetadataRepository, fileid uuid.UUID, storedRevision **Revision) {
	t.Helper()
	mockFileMetadataRepository.EXPECT().StoreRevision(mock.Anything, fileid, mock.IsType(Revision{})).
		Run(func(ctx context.Context, existingFileId uuid.UUID, revision Revision) { *storedRevision = &revision }).
		Return(nil).
		Times(1)
}

func setupFileMetadataRepositoryMockStoreFileMetadataReturnsError(t *testing.T, mockFileMetadataRepository *MockFileMetadataRepository, err error) *MockFileMetadataRepository {
	t.Helper()
	mockFileMetadataRepository.EXPECT().StoreFileMetadata(mock.Anything, mock.IsType(FileMetadata{})).
		Return(err).
		Times(1)

	return mockFileMetadataRepository
}

func setupFileMetadataRepositoryMockStoreRevisionReturnsError(t *testing.T, mockFileMetadataRepository *MockFileMetadataRepository, fileId uuid.UUID, err error) {
	t.Helper()
	mockFileMetadataRepository.EXPECT().StoreRevision(mock.Anything, fileId, mock.IsType(Revision{})).
		Return(err).
		Times(1)
}

func setupFileMetadataRepositoryToFetchMetadata(t *testing.T, mockFileMetadataRepository *MockFileMetadataRepository, fileId uuid.UUID, fileMetadata FileMetadata) {
	t.Helper()
	mockFileMetadataRepository.EXPECT().FetchFileMetadata(mock.Anything, fileId).Return(fileMetadata, nil).Times(1)
	mockFileMetadataRepository.EXPECT().NotFoundError().Return(errors.New("not expected error")).Times(1)
}
