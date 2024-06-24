package api_platform_client

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

// # Request `accessToken` from API Platform.
//
// This function is used to request a `accessToken` from API Platform.
// Note that every token is expired after a day.
func (client *ApiPlatformClient) RequestAccessToken() (string, error) {
	endpointAccessToken := fmt.Sprintf("%s/tsmpaa/oauth/token", client.EndpointURL)

	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	header.Set("Authorization", "Basic "+basicAuthCredentialCrafter(client.ClientID, client.ClientTokenBlock))

	payload := url.Values{}
	payload.Set("grant_type", "client_credentials")

	resp, err := http.PostForm(endpointAccessToken, payload)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return "", ErrInvalidCredential
		}
		return "", ErrApiPlatformGeneralError
	}

	tokenResponse := new(ApiPlatformTokenResponse)
	if err := json.NewDecoder(resp.Body).Decode(tokenResponse); err != nil {
		return "", err
	}

	client.accessToken = tokenResponse.AccessToken
	return client.accessToken, nil
}

func (client *ApiPlatformClient) SignPayload(data []byte) string {
	signBody := append([]byte(client.signBlock), data...)
	signature := sha256.Sum256(signBody)
	return fmt.Sprintf("%x", signature)
}

// # Payload Crafter for HTTP Basic Auth
//
// This function is used to craft the payload for HTTP Basic Auth.
func basicAuthCredentialCrafter(clientID, clientTokenBlock string) string {
	credential := fmt.Sprintf("%s:%s", clientID, clientTokenBlock)
	return base64.StdEncoding.EncodeToString([]byte(credential))
}
