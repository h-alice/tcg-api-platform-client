package api_platform_client

type ApiPlatformTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expires     int    `json:"expires_in"`
	Scope       string `json:"scope"`
	Node        string `json:"node"`
	Jti         string `json:"jti"` // Unknown
}

type ApiPlatformSignBlockResHeaderResponse struct {
	ReturnCode    string `json:"rtnCode"`
	ReturnMessage string `json:"rtnMsg"`
}

type ApiPlatformSignBlockResResponse struct {
	SignBlock string `json:"signBlock"`
}

type ApiPlatformSignBlockResponse struct {
	ResponseHeader    ApiPlatformSignBlockResHeaderResponse `json:"ResHeader"`
	ResponseSignBlock ApiPlatformSignBlockResResponse       `json:"Res_getSignBlock"`
}
