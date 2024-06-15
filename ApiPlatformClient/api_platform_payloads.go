package api_platform_client

type ApiPlatformTokenResponse struct {
	AccessToken string
	TokenType   string
	Expires     int
	Scope       string
	Node        string
	Jti         string // Unknown
}
