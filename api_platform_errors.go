package api_platform_client

import (
	"errors"
)

var (
	ErrInvalidCredential       = errors.New("invalid client credentials, check `ClientID` and `ClientTokenBlock`")
	ErrClientUnauthorized      = errors.New("client unauthorized")
	ErrApiPlatformGeneralError = errors.New("api platform general error")
)
