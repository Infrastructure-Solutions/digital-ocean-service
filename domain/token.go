package domain

// DOToken is the response of the OAUTH for Digital Ocean
type DOToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	Info         Info   `json:"info"`
}

// Info is the user information in a OAUTH requests from Digital Ocean
type Info struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	UUID  string `json:"uuid"`
}
