package errs

import "errors"

var ErrAlreadyExists = errors.New("already exists error")
var ErrNotExists = errors.New("not exists error")
var ErrOrderLoadAnotherUser = errors.New("the order number has already been uploaded by another user")
var ErrOrderAlreadyLoadThisUser = errors.New("the order number has already been uploaded by this user")
var ErrOrderNum = errors.New("invalid order number format")
var ErrNotEnoughFunds = errors.New("there are not enough funds in the account")
var ErrInternalServerError = errors.New("internal server error")
