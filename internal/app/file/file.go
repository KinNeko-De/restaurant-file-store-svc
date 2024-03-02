package file

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

const sniffSize = 512 // defined by the net/http package

type FileMetadata struct {
	Id        uuid.UUID
	Revisions []Revision
}

func (f *FileMetadata) AddRevision(revision Revision) {
	f.Revisions = append(f.Revisions, revision)
}

func (f *FileMetadata) LatestRevision() Revision {
	return f.Revisions[len(f.Revisions)-1]
}

func (f *FileMetadata) FirstRevision() Revision {
	return f.Revisions[0]
}

func (f *FileMetadata) CreatedAt() time.Time {
	return f.FirstRevision().CreatedAt
}

func (f *FileMetadata) LastUpdatedAt() time.Time {
	return f.LatestRevision().CreatedAt
}

type Revision struct {
	Id        uuid.UUID
	Extension string
	MediaType string
	Size      uint64
	CreatedAt time.Time
}

func newFileMetadata(fileId uuid.UUID, initialRevision Revision) FileMetadata {
	return FileMetadata{
		Id:        fileId,
		Revisions: []Revision{initialRevision},
	}
}

func newRevision(fileName string, fileSize uint64, sniff []byte) Revision {
	createdAt := time.Now().UTC()

	return Revision{
		Id:        uuid.New(),
		Extension: filepath.Ext(fileName),
		MediaType: http.DetectContentType(sniff),
		Size:      fileSize,
		CreatedAt: createdAt,
	}
}
