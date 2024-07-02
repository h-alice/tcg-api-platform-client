package api_platform_client

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	json_ "encoding/json"
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

	// The endpoint to request the `accessToken`.
	endpointAccessToken := fmt.Sprintf("%s/tsmpaa/oauth/token", client.EndpointURL)

	payload := url.Values{}                         // The payload to request the `accessToken`.
	payload.Set("grant_type", "client_credentials") // Add proper grant type.

	// Create a new HTTP request.
	req, err := http.NewRequest("POST", endpointAccessToken, bytes.NewBufferString(payload.Encode()))
	if err != nil {
		return "", err
	}

	// Set the proper headers.
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set the Authorization header with encoded `clientID` and `clientTokenBlock`.
	req.Header.Set("Authorization", "Basic "+basicAuthCredentialCrafter(client.ClientID, client.ClientTokenBlock))

	// Create a new HTTP client.
	http_client := &http.Client{}

	// Send the request.
	resp, err := http_client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() // Close the response body after the function ends.

	// Check if the response status code is not OK.
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized { // 401 Unauthorized
			return "", ErrInvalidCredential // Return invalid credential error.
		}
		return "", ErrApiPlatformGeneralError // TODO: Add more specific error.
	}

	// Decode the response body to `ApiPlatformTokenResponse` struct.
	tokenResponse := new(ApiPlatformTokenResponse)
	if err := json_.NewDecoder(resp.Body).Decode(tokenResponse); err != nil {
		return "", err
	}

	// Set the `accessToken` to the client.
	client.accessToken = tokenResponse.AccessToken

	return client.accessToken, nil // Return the `accessToken`.
}

// # Request `signBlock` from API Platform.
//
// This function is used to request a `signBlock` from API Platform.
// Every payload body must be signed with the `signBlock` before sending it to the API Platform.
func (client *ApiPlatformClient) RequestSignBlock() (string, error) {

	// Initial check if the client is not authorized.
	if client.accessToken == "" {
		return "", ErrClientUnauthorized // Return client unauthorized error.
	}

	// The endpoint to request the `signBlock`.
	endpointSB := fmt.Sprintf("%s/tsmpaa/getSignBlock", client.EndpointURL)

	// Create a new HTTP request.
	req, err := http.NewRequest("GET", endpointSB, nil)
	if err != nil {
		return "", err
	}

	// Set the Authorization header with the `accessToken`.
	req.Header.Set("Authorization", "Bearer "+client.accessToken)

	// Create a new HTTP client.
	clientHTTP := &http.Client{}

	// Send the request.
	resp, err := clientHTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() // Close the response body after the function ends.

	// Check if the response status code is not OK.
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized { // 401 Unauthorized
			return "", ErrInvalidCredential // Return invalid credential error.
		}
		return "", ErrApiPlatformGeneralError // TODO: Add more specific error.
	}

	// Decode the response body to `ApiPlatformSignBlockResponse` struct.
	sbResponse := new(ApiPlatformSignBlockResponse)
	if err := json_.NewDecoder(resp.Body).Decode(sbResponse); err != nil {
		return "", err
	}

	// Set the `signBlock` to the client.
	client.signBlock = sbResponse.ResponseSignBlock.SignBlock
	return client.signBlock, nil // Return the `signBlock`.
}

// # Sign Payload
//
// This function is used to sign the payload with the `signBlock`.
// This function is expected to be used after client authorized and `signBlock` is generated.
func (client *ApiPlatformClient) SignPayload(data []byte) string {
	signBody := append([]byte(client.signBlock), data...)
	signature := sha256.Sum256(signBody)
	return fmt.Sprintf("%x", signature)
}

// # Send Request
//
// This function is used to send a request via the API Platform.
// This function designed to imitate the `requests.Request` function in Python.
func (client *ApiPlatformClient) SendRequest(endpoint, method string, headers map[string]string, jsonPayload, data interface{}) (*http.Response, error) {
	if client.accessToken == "" { // Initial check if the client is not authorized.
		return nil, ErrClientUnauthorized // Return client unauthorized error.
	}

	if client.signBlock == "" { // Initial check if the `signBlock` is not generated.
		return nil, ErrClientUnauthorized // Return client unauthorized error.
	}

	// Payload cfrating.
	jsonData, err := json_.Marshal(jsonPayload) // Marshal the JSON payload.
	if err != nil {
		return nil, err
	}

	// Sign the payload.
	signature := client.SignPayload(jsonData)

	// Create a new HTTP request.
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Set the Authorization header with the `accessToken`.
	req.Header.Set("Authorization", "Bearer "+client.accessToken)
	req.Header.Set("SignCode", signature)

	// Add additional header fields.
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Create a new HTTP client.
	clientHTTP := &http.Client{}
	return clientHTTP.Do(req) // Send the request.
}

// # New API Platform Client
//
// This function is used to create a new API Platform client.
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
