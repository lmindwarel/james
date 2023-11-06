package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// UUID is the type for UUIDs
type UUID string

type BaseModel struct {
	ID UUID `json:"id" bson:"_id"`
	// DateCreated is the date when the object was created
	DateCreated time.Time `json:"dateCreated,omitempty" bson:"dateCreated"`
	// DateModified is the date when the object was updated
	DateUpdated time.Time `json:"dateUpdated,omitempty" bson:"dateUpdated"`
}

func (m BaseModel) Update() {
	m.DateUpdated = time.Now()
}

func NewBaseModel() BaseModel {
	now := time.Now()
	return BaseModel{
		ID:          NewUUID(),
		DateCreated: now,
		DateUpdated: now,
	}
}

func EmptyUUID(id UUID) bool {
	return string(id) == ""
}

// NewUUID return new UUID
func NewUUID() UUID {
	return UUID(uuid.NewV4().String())
}
