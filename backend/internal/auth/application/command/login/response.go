package login

type Response struct {
	AccessToken string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}