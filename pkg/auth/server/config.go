package server

type Config struct {
	IssuerURL      string
	OfflineAsScope bool
	RootToken      *string
}
