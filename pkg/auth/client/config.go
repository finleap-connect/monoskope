package client

type Config struct {
	IssuerURL      string
	OfflineAsScope bool
	Nonce          string
	ClientSecret   string
	RedirectURL    string
}
