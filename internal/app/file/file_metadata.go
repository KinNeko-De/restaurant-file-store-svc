package file

import (
	"time"

	"github.com/google/uuid"
)

type FileMetadata struct {
	FileId     uuid.UUID `bson:"_id"`
	RevisionId uuid.UUID `bson:"revision"`
	Extension  string    `bson:"extension"`
	Size       uint64    `bson:"size"`
	MediaType  string    `bson:"mediaType"`
	CreatedAt  time.Time `bson:"createdAt"`
}
