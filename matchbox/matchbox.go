package matchbox

// Config configures a matchbox client.
type Config struct {
	// gRPC API endpoint
	Endpoint string
	// PEM encoded TLS CA and client credentials
	CA         []byte
	ClientCert []byte
	ClientKey  []byte
}
