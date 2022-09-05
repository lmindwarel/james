package models

import (
	"fmt"
	"time"
)

// Rating hold informations about user rating on element (Quizz, etc.)
type Rating struct {
	// Value is the rating value (Generally between 0 and 5 stars)
	Value float32 `json:"value,omitempty" bson:"value"`
	// Raters is a map containing each user rating value
	Raters map[UUID]float32 `json:"raters,omitempty" bson:"raters"`
}

const (
	DefaultPaginationLimit = 100
	MaxPaginationLimit     = 200
)

type Pagination struct {
	Offset int
	Limit  int
	Sort   []string // All column that was sorted
}

var DefaultPaginationFilters = Pagination{
	Limit: DefaultPaginationLimit,
}

type PaginatedResults struct {
	Pagination
	Total int64
}

type DurationMs time.Duration

func (d DurationMs) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf("%d", time.Duration(d).Milliseconds())), nil
}
