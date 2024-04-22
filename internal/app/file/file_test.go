//go:build unit

package file

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	fixture "github.com/kinneko-de/restaurant-file-store-svc/test/testing/file"
)

func TestFileMetadata_AddRevision(t *testing.T) {
	fileId := uuid.New()
	revision := Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain",
		Size:      1024,
		CreatedAt: time.Now().UTC(),
	}

	fileMetadata := FileMetadata{
		Id:        fileId,
		Revisions: []Revision{},
	}

	fileMetadata.AddRevision(revision)

	assert.Equal(t, 1, len(fileMetadata.Revisions))
	assert.Equal(t, revision, fileMetadata.Revisions[0])
}

func TestFileMetadata_LatestRevision_TwoRevisions(t *testing.T) {
	fileId := uuid.New()
	revision1 := Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain; charset=utf-8",
		Size:      1024,
		CreatedAt: time.Now().UTC(),
	}
	expectedRevision := Revision{
		Id:        uuid.New(),
		Extension: ".pdf",
		MediaType: "application/pdf",
		Size:      2048,
		CreatedAt: time.Now().UTC(),
	}

	fileMetadata := &FileMetadata{
		Id:        fileId,
		Revisions: []Revision{revision1, expectedRevision},
	}

	latestRevision := fileMetadata.LatestRevision()

	assert.Equal(t, expectedRevision, latestRevision)
}

func TestFileMetadata_firstRevision_TwoRevisions(t *testing.T) {
	fileId := uuid.New()
	expectedRevision := Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain; charset=utf-8",
		Size:      1024,
		CreatedAt: time.Now().UTC(),
	}
	revision2 := Revision{
		Id:        uuid.New(),
		Extension: ".pdf",
		MediaType: "application/pdf",
		Size:      2048,
		CreatedAt: time.Now().UTC(),
	}

	fileMetadata := &FileMetadata{
		Id:        fileId,
		Revisions: []Revision{expectedRevision, revision2},
	}

	FirstRevision := fileMetadata.firstRevision()

	assert.Equal(t, expectedRevision, FirstRevision)
}

func TestFileMetadata_LastUpdatedAt_TwoRevisions(t *testing.T) {
	fileId := uuid.New()
	revision1 := Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain; charset=utf-8",
		Size:      1024,
		CreatedAt: time.Now().UTC().Add(-time.Hour),
	}
	revision2 := Revision{
		Id:        uuid.New(),
		Extension: ".pdf",
		MediaType: "application/pdf",
		Size:      2048,
		CreatedAt: time.Now().UTC(),
	}

	fileMetadata := &FileMetadata{
		Id:        fileId,
		Revisions: []Revision{revision1, revision2},
	}

	lastUpdatedAt := fileMetadata.LastUpdatedAt()

	assert.Equal(t, revision2.CreatedAt, lastUpdatedAt)
}

func TestFileMetadata_CreatedAt_TwoRevisions(t *testing.T) {
	fileId := uuid.New()
	revision1 := Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain; charset=utf-8",
		Size:      1024,
		CreatedAt: time.Now().UTC().Add(-time.Hour),
	}
	revision2 := Revision{
		Id:        uuid.New(),
		Extension: ".pdf",
		MediaType: "application/pdf",
		Size:      2048,
		CreatedAt: time.Now().UTC(),
	}

	fileMetadata := &FileMetadata{
		Id:        fileId,
		Revisions: []Revision{revision1, revision2},
	}

	createdAt := fileMetadata.CreatedAt()

	assert.Equal(t, revision1.CreatedAt, createdAt)
}

func TestNewFileMetadata(t *testing.T) {
	fileId := uuid.New()
	revision := Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain; charset=utf-8",
		Size:      1024,
		CreatedAt: time.Now().UTC(),
	}

	fileMetadata := newFileMetadata(fileId, revision)

	assert.Equal(t, fileId, fileMetadata.Id)
	assert.Equal(t, []Revision{revision}, fileMetadata.Revisions)
}

func TestNewRevision(t *testing.T) {
	revisionId := uuid.New()
	extension := ".txt"
	fileSize := uint64(4)
	sniff := fixture.TextFile()
	expectedExtension := ".txt"
	exectedMediaType := "text/plain; charset=utf-8"

	revision := newRevision(revisionId, extension, fileSize, sniff)

	assert.NotEqual(t, uuid.Nil, revision.Id)
	assert.Equal(t, expectedExtension, revision.Extension)
	assert.Equal(t, exectedMediaType, revision.MediaType)
	assert.Equal(t, fileSize, revision.Size)
	assert.WithinDuration(t, time.Now().UTC(), revision.CreatedAt, time.Second)
}

func TestFileMetadata_GetRevision(t *testing.T) {
	fileId := uuid.New()
	expectedRevision := Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain; charset=utf-8",
		Size:      1024,
		CreatedAt: time.Now().UTC().Add(-time.Hour),
	}
	latestRevision := Revision{
		Id:        uuid.New(),
		Extension: ".pdf",
		MediaType: "application/pdf",
		Size:      2048,
		CreatedAt: time.Now().UTC(),
	}

	sut := &FileMetadata{
		Id:        fileId,
		Revisions: []Revision{expectedRevision, latestRevision},
	}

	actualRevision, err := sut.GetRevision(expectedRevision.Id)

	assert.NoError(t, err)
	assert.Equal(t, expectedRevision, actualRevision)
}

func TestFileMetadata_GetRevision_NotFound(t *testing.T) {
	requestedRevision := uuid.New()
	existingRevision := uuid.New()

	latestRevision := Revision{
		Id:        existingRevision,
		Extension: ".pdf",
		MediaType: "application/pdf",
		Size:      2048,
		CreatedAt: time.Now().UTC(),
	}

	sut := &FileMetadata{
		Id:        uuid.New(),
		Revisions: []Revision{latestRevision},
	}

	actualRevision, err := sut.GetRevision(requestedRevision)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), requestedRevision.String())
	assert.Equal(t, Revision{}, actualRevision)
}
