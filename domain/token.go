package domain

type DOToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	Info         Info   `json:"info"`
}

type Info struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	UUID  string `json:"uuid"`
}
