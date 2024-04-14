package file

import (
	"fmt"
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

func (f *FileMetadata) GetRevision(revisionId uuid.UUID) (Revision, error) {
	for _, revision := range f.Revisions {
		if revision.Id == revisionId {
			return revision, nil
		}
	}
	return Revision{}, fmt.Errorf("revision '%s' not found", revisionId)
}

func (f *FileMetadata) firstRevision() Revision {
	return f.Revisions[0]
}

func (f *FileMetadata) CreatedAt() time.Time {
	return f.firstRevision().CreatedAt
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

func newRevision(revisionId uuid.UUID, fileName string, fileSize uint64, sniff []byte) Revision {
	createdAt := time.Now().UTC()

	return Revision{
		Id:        revisionId,
		Extension: filepath.Ext(fileName),
		MediaType: http.DetectContentType(sniff),
		Size:      fileSize,
		CreatedAt: createdAt,
	}
}
