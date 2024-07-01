# Taipei City Government API Management Platform Client

This project provides a Go client for sending requests via the "API management platform" of the Taipei City Government.

Currently we only support sending JSON requests for now. Sending other types of requests (e.g., form data) will be supported in the future.

## ✨ Features

- Request access tokens / sign blocks in simple way.
- Send signed requests to the Taipei City Government API.

## 📦 Installation

To use this client, you need to install the package. You can do this using `go get`:

```bash
go get github.com/h-alice/tcg-api-platform-client
```

## 🚀 Usage

Below is a sample code demonstrating how to use the client to interact with the "TaipeiON Chat" service via the API management platform. Make sure to replace all secret credentials with your actual credentials and adjust the API endpoint and format to fit your requirements.

```go
package main

import (
	"fmt"
	"io"

	api_client "github.com/h-alice/tcg-api-platform-client"
)

func main() {
	// Example usage, note that you may want to change the endpoint if you switched from test server to production server.
	client := api_client.NewApiPlatformClient("https://apimtest.gov.taipei", "[your_client_id]", "[your_client_token_block]")

	// Request Access Token
	accessToken, err := client.RequestAccessToken()
	if err != nil {
		fmt.Println("Error requesting access token:", err)
		return
	}

	fmt.Println("Access Token:", accessToken)

	// Request Sign Block
	signBlock, err := client.RequestSignBlock()
	if err != nil {
		fmt.Println("Error requesting sign block:", err)
		return
	}
	fmt.Println("Sign Block:", signBlock)


    // Add custom headers and payload. In this example, the `backAuth` field is for TaipeiON Chat service authentication.
	header := map[string]string{
		"Content-Type": "application/json",
		"backAuth":     "[back_auth_token_of_taipeion]",
	}

    // Broadcast a message to all users subscribed to channel.
	payload := map[string]interface{}{
		"ask": "broadcastMessage",
		"message": map[string]interface{}{
			"type": "text",
			"text": " Hello World from API platform!",
		},
	}

    // Send Request to TaipeiON Chat service via API management platform.
	resp, err := client.SendRequest("https://apimtest.gov.taipei/tsmpc/m-taipeion/MessageFeedService", "POST", header, payload, nil)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response:", string(body)) // Print the response body.
}
```

## 📖 Methods

### NewApiPlatformClient

Creates a new `ApiPlatformClient` instance.

**Parameters:**

- `endpointURL` (string): The endpoint URL of the API management platform.
- `clientID` (string): The client ID, you can check the registration information in the API management platform.
- `clientTokenBlock` (string): The client token block, you can check the registration information in the API management platform.

**Returns:**

- `*ApiPlatformClient`: A new instance of `ApiPlatformClient`.

### RequestAccessToken

Requests an access token from the API.

**Returns:**

- `string`: The access token.
- `error`: An error if the request fails.

### RequestSignBlock

Requests a sign block from the API.

**Returns:**

- `string`: The sign block.
- `error`: An error if the request fails.

### SendRequest

Sends a signed request to the API.

**Parameters:**

- `endpoint` (string): The API endpoint.
- `method` (string): The HTTP method (e.g., "POST").
- `headers` (map[string]string): The request headers.
- `jsonPayload` (interface{}): The JSON payload.
- `data` (interface{}): Additional data.

**Returns:**

- `*http.Response`: The response from the API.
- `error`: An error if the request fails.

## 🤝 Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## 📜 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
