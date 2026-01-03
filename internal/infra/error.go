package infra

import "errors"

var ErrErrorResponse = errors.New("error response from external api")
var ErrTemporaryUnavaillalble = errors.New("temporary unavailable error")
var ErrUnexpected = errors.New("unexpected error")
