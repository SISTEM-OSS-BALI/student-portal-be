package auth

import "errors"

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidResetToken = errors.New("invalid reset token")
