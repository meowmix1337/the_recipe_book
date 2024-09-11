package endpoint

type JWTResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
