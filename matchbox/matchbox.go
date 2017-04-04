package matchbox

// MatchboxConfig configures a matchbox client with PEM encoded TLS credentials.
type MatchboxConfig struct {
	Endpoint   string
	CA         []byte
	ClientCert []byte
	ClientKey  []byte
}
