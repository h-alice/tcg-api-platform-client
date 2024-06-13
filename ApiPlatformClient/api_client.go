package api_platform_client

import (
	"encoding/base64"
	"fmt"
)

// # API Platform Client Struct
//
// This struct is used to store the API platform client information.
type ApiPlatformClient struct {
	EndpointURL      string // The URL of the API Platform. Can be either test or production.
	ClientID         string
	ClientTokenBlock string
	accessToken      string // `accessToken` is generated during runtime.
	signBlock        string // `signBlock` is generated during runtime.
}

// # Payload Crafter for HTTP Basic Auth
//
// This function is used to craft the payload for HTTP Basic Auth.
func basicAuthCredentialCrafter(clientID, clientTokenBlock string) string {
	credential := fmt.Sprintf("%s:%s", clientID, clientTokenBlock)
	return base64.StdEncoding.EncodeToString([]byte(credential))
}
