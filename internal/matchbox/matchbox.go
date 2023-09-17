package matchbox

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"time"

	matchbox "github.com/poseidon/matchbox/matchbox/client"
)

var (
	defaultTimeout = 25 * time.Second
)

// Config configures a matchbox client.
type Config struct {
	// gRPC API endpoint
	Endpoint string
	// PEM encoded TLS CA and client credentials
	CA         []byte
	ClientCert []byte
	ClientKey  []byte
}

// NewMatchboxClient returns a new matchbox.Client.
func NewMatchboxClient(config *MatchboxProviderModel) (*matchbox.Client, error) {
	tlscfg, err := tlsConfig(
		[]byte(config.CA.ValueString()),
		[]byte(config.ClientCert.ValueString()),
		[]byte(config.ClientKey.ValueString()),
	)
	if err != nil {
		return nil, err
	}
	return matchbox.New(&matchbox.Config{
		Endpoints:   []string{config.Endpoint.ValueString()},
		DialTimeout: defaultTimeout,
		TLS:         tlscfg,
	})
}

// tlsConfig returns a matchbox client TLS.Config.
// TODO: Update matchbox TLSInfo.ClientConfig to replace this.
func tlsConfig(ca, clientCert, clientKey []byte) (*tls.Config, error) {
	// certificate authority for verifying the server
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(ca)
	if !ok {
		return nil, errors.New("no PEM certificates were parsed")
	}

	// client certificate for authentication
	cert, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		// CA bundle the client should trust when verifying the server
		RootCAs: pool,
		// Client certificate to authenticate to the server
		Certificates: []tls.Certificate{cert},
	}, nil
}
