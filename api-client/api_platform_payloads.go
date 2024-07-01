package api_platform_client

// # Response structure for the API Platform Token request
type ApiPlatformTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expires     int    `json:"expires_in"`
	Scope       string `json:"scope"`
	Node        string `json:"node"`
	Jti         string `json:"jti"` // Unknown
}

type apiPlatformSignBlockResHeaderResponse struct {
	ReturnCode    string `json:"rtnCode"`
	ReturnMessage string `json:"rtnMsg"`
}

type apiPlatformSignBlockResResponse struct {
	SignBlock string `json:"signBlock"`
}

// # Response structure for the API Platform SignBlock request
type ApiPlatformSignBlockResponse struct {
	ResponseHeader    apiPlatformSignBlockResHeaderResponse `json:"ResHeader"`
	ResponseSignBlock apiPlatformSignBlockResResponse       `json:"Res_getSignBlock"`
}
