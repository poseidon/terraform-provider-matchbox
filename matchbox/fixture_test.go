package matchbox

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/coreos/matchbox/matchbox/rpc"
	"github.com/coreos/matchbox/matchbox/server"
	"github.com/coreos/matchbox/matchbox/storage"
	"github.com/coreos/matchbox/matchbox/tlsutil"
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
	Store     storage.Store
	Servers   []*grpc.Server
	Listeners []net.Listener
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

func NewFixtureServer(clientTLS *TLSContents, serverTLS *tlsutil.TLSInfo, s storage.Store, replicas int) *FixtureServer {
	listeners := []net.Listener{}
	// Address close (i.e. release) is effectively asynchronous. Test server
	// instances should reserve a random address upfront.
	for i := 0; i < replicas; i++ {
		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			panic(fmt.Errorf("failed to start listening: %v", err))
		}

		listeners = append(listeners, lis)
	}

	return &FixtureServer{
		Store:     s,
		Servers:   []*grpc.Server{},
		Listeners: listeners,
		ServerTLS: serverTLS,
		ClientTLS: clientTLS,
	}
}

func (s *FixtureServer) Start() {
	cfg, err := s.ServerTLS.ServerConfig()
	if err != nil {
		panic(fmt.Errorf("Invalid TLS credentials: %v", err))
	}

	srv := server.NewServer(&server.Config{Store: s.Store})
	for _, listener := range s.Listeners {
		server := rpc.NewServer(srv, cfg)
		s.Servers = append(s.Servers, server)
		go server.Serve(listener)
	}
}

func (s *FixtureServer) Stop() {
	if len(s.Servers) > 0 {
		for _, server := range s.Servers {
			server.Stop()
		}
		s.Servers = []*grpc.Server{}
	}
}

func (s *FixtureServer) AddProviderConfig(hcl string) string {
	provider := `
		provider "matchbox" {
			%s
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
	endpoints := ""
	if len(s.Listeners) == 1 {
		endpoints = fmt.Sprintf(`endpoint = "%s"`, s.Listeners[0].Addr().String())
	} else {
		addresses := ""
		for _, listener := range s.Listeners {
			addresses += fmt.Sprintf(`"%s",`, listener.Addr().String())
		}
		addresses = strings.Trim(addresses, ",")
		endpoints = fmt.Sprintf(`endpoints = [%s]`, addresses)
	}
	return fmt.Sprintf(provider,
		endpoints,
		s.ClientTLS.Cert,
		s.ClientTLS.Key,
		s.ClientTLS.CA,
		hcl)
}

// mustFile wraps a call to ioutil.ReadFile and panics if the error is non-nil.
func mustReadFile(filename string) []byte {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return contents
}
