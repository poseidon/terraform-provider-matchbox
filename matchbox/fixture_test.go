package matchbox

import (
	"fmt"
	"net"
	"os"

	"github.com/poseidon/matchbox/matchbox/rpc"
	"github.com/poseidon/matchbox/matchbox/server"
	"github.com/poseidon/matchbox/matchbox/storage"
	"github.com/poseidon/matchbox/matchbox/tlsutil"
	"google.golang.org/grpc"
)

var (
	fakeClientCert = mustReadFile("testdata/client.crt")
	fakeClientKey  = mustReadFile("testdata/client.key")
	fakeCACert     = mustReadFile("testdata/ca.crt")
	clientTLSInfo  = &TLSContents{
		Cert: fakeClientCert,
		Key:  fakeClientKey,
		CA:   fakeCACert,
	}
	serverTLSInfo = &tlsutil.TLSInfo{
		CAFile:   "testdata/ca.crt",
		CertFile: "testdata/server.crt",
		KeyFile:  "testdata/server.key",
	}
)

type FixtureServer struct {
	Store    storage.Store
	Server   *grpc.Server
	Listener net.Listener
	// TLS server certificates (files) which will be used
	ServerTLS *tlsutil.TLSInfo
	// TLS client credentials which should be used
	ClientTLS *TLSContents
}

// TODO: Merge into matchbox tlsutil TLSInfo to allow raw contents.
type TLSContents struct {
	Cert []byte
	Key  []byte
	CA   []byte
}

func NewFixtureServer(clientTLS *TLSContents, serverTLS *tlsutil.TLSInfo, s storage.Store) *FixtureServer {
	// Address close (i.e. release) is effectively asynchronous. Test server
	// instances should reserve a random address upfront.
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(fmt.Errorf("failed to start listening: %v", err))
	}
	return &FixtureServer{
		Store:     s,
		Listener:  lis,
		ServerTLS: serverTLS,
		ClientTLS: clientTLS,
	}
}

func (s *FixtureServer) Start() error {
	cfg, err := s.ServerTLS.ServerConfig()
	if err != nil {
		return fmt.Errorf("Invalid TLS credentials: %v", err)
	}

	srv := server.NewServer(&server.Config{Store: s.Store})
	s.Server = rpc.NewServer(srv, cfg)
	return s.Server.Serve(s.Listener)
}

func (s *FixtureServer) Stop() {
	if s.Server != nil {
		s.Server.Stop()
	}
}

func (s *FixtureServer) AddProviderConfig(hcl string) string {
	provider := `
		provider "matchbox" {
			endpoint = "%s"
			client_cert = <<CERT
%s
CERT
			client_key = <<KEY
%s
KEY
			ca         = <<CA
%s
CA
		}

		%s
		`
	return fmt.Sprintf(provider,
		s.Listener.Addr().String(),
		s.ClientTLS.Cert,
		s.ClientTLS.Key,
		s.ClientTLS.CA,
		hcl)
}

// mustFile wraps a call to ioutil.ReadFile and panics if the error is non-nil.
func mustReadFile(filename string) []byte {
	contents, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return contents
}
