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

func TestFileMetadata_FirstRevision_TwoRevisions(t *testing.T) {
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

	FirstRevision := fileMetadata.FirstRevision()

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
	fileName := "test.txt"
	fileSize := uint64(4)
	sniff := fixture.TextFile()
	expectedExtension := ".txt"
	exectedMediaType := "text/plain; charset=utf-8"

	revision := newRevision(fileName, fileSize, sniff)

	assert.NotEqual(t, uuid.Nil, revision.Id)
	assert.Equal(t, expectedExtension, revision.Extension)
	assert.Equal(t, exectedMediaType, revision.MediaType)
	assert.Equal(t, fileSize, revision.Size)
	assert.WithinDuration(t, time.Now().UTC(), revision.CreatedAt, time.Second)
}
