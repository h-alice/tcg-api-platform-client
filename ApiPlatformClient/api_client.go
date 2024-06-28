package api_platform_client

import (
	"bytes"
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

func (client *ApiPlatformClient) RequestSignBlock() (string, error) {
	if client.accessToken == "" {
		return "", ErrClientUnauthorized
	}

	endpointSB := fmt.Sprintf("%s/tsmpaa/getSignBlock", client.EndpointURL)

	req, err := http.NewRequest("GET", endpointSB, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+client.accessToken)

	clientHTTP := &http.Client{}
	resp, err := clientHTTP.Do(req)
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

	sbResponse := new(ApiPlatformSignBlockResponse)
	if err := json.NewDecoder(resp.Body).Decode(sbResponse); err != nil {
		return "", err
	}

	client.signBlock = sbResponse.ResponseSignBlock.SignBlock
	return client.signBlock, nil
}

func (client *ApiPlatformClient) SignPayload(data []byte) string {
	signBody := append([]byte(client.signBlock), data...)
	signature := sha256.Sum256(signBody)
	return fmt.Sprintf("%x", signature)
}

func (client *ApiPlatformClient) SendRequest(endpoint, method string, headers map[string]string, jsonPayload, data interface{}) (*http.Response, error) {
	if client.accessToken == "" {
		return nil, ErrClientUnauthorized
	}

	if client.signBlock == "" {
		return nil, ErrClientUnauthorized
	}

	jsonData, err := json.Marshal(jsonPayload)
	if err != nil {
		return nil, err
	}

	signature := client.SignPayload(jsonData)

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+client.accessToken)
	req.Header.Set("SignCode", signature)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	clientHTTP := &http.Client{}
	return clientHTTP.Do(req)
}

func NewApiPlatformClient(endpointURL, clientID, clientTokenBlock string) *ApiPlatformClient {
	return &ApiPlatformClient{
		EndpointURL:      endpointURL,
		ClientID:         clientID,
		ClientTokenBlock: clientTokenBlock,
	}
}

// # Payload Crafter for HTTP Basic Auth
//
// This function is used to craft the payload for HTTP Basic Auth.
func basicAuthCredentialCrafter(clientID, clientTokenBlock string) string {
	credential := fmt.Sprintf("%s:%s", clientID, clientTokenBlock)
	return base64.StdEncoding.EncodeToString([]byte(credential))
}
