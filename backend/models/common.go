package models

import (
	"fmt"
	"time"
)

// Rating hold informations about user rating on element (Quizz, etc.
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

type CredentialPatch struct {
	User     *string
	Password *string
}
