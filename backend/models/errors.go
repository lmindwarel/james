package models

import "errors"

var (
	ErrBadParameter         = errors.New("bad parameter")
	ErrAuthenticationFailed = errors.New("authentication failed")
)
