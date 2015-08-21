package domain

type User struct {
	ID    int64
	Token DOToken
}

type Key struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	PublicKey   string `json:"public_key,omitempty"`
}
