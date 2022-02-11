package common

import "errors"

var (
	ErrEmptyAuthHeader        = errors.New("authorization header is empty")
	ErrMalformedToken         = errors.New("malformed Token")
	ErrBadRouting             = errors.New("bad routing")
	ErrKIDNotFound            = errors.New("kid header not found")
	ErrUnableToParsePublicKey = errors.New("could not parse public key")
	ErrUnexpectedTokenVersion = errors.New("unexpected token version")
	ErrJwtTokenInvalid        = errors.New("invalid JWT token")
)
