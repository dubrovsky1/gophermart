package errs

import "errors"

var ErrAlreadyExists = errors.New("already exists error")
var ErrNotExists = errors.New("not exists error")
