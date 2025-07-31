package storageerror

import "errors"

var (
	ErrSubAlredyExist = errors.New("sub alredy exist")
	ErrEmptyStorage   = errors.New("sub storage is empty")
	ErrSubNoFound     = errors.New("sub not found")
)
