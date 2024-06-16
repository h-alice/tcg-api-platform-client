package api_platform_client

type ApiPlatformTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expires     int    `json:"expires_in"`
	Scope       string `json:"scope"`
	Node        string `json:"node"`
	Jti         string `json:"jti"` // Unknown
}
