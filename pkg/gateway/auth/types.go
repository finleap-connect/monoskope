package auth

type Config struct {
	IssuerURL      string
	OfflineAsScope bool
}

type ExtraClaims struct {
	Email         string   `json:"email"`
	EmailVerified bool     `json:"email_verified"`
	Groups        []string `json:"groups"`
}
