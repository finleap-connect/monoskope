package client

type Config struct {
	IssuerURL      string
	OfflineAsScope bool
	Nonce          string
	ClientId       string
	ClientSecret   string
	RedirectURI    string
}
