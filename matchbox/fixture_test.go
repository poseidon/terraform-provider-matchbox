package matchbox

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/coreos/matchbox/matchbox/rpc"
	"github.com/coreos/matchbox/matchbox/server"
	"github.com/coreos/matchbox/matchbox/storage"
)

type FixtureServer struct {
	Store    storage.Store
	Server   *grpc.Server
	Listener net.Listener
}

func NewFixtureServer(s storage.Store) *FixtureServer {
	return &FixtureServer{
		Store: s,
	}
}

func (s *FixtureServer) Start() {
	var err error
	s.Listener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(fmt.Errorf("failed to start listener: %v", err))
	}

	cfg, err := testServerTLSConfig()
	if err != nil {
		panic(fmt.Errorf("Invalid TLS credentials: %v", err))
	}

	srv := server.NewServer(&server.Config{s.Store})
	s.Server = rpc.NewServer(srv, cfg)
	go s.Server.Serve(s.Listener)
}

func (s *FixtureServer) Stop() {
	s.Server.Stop()
}

func (s *FixtureServer) AddProviderConfig(hcl string) string {
	return fmt.Sprintf(`
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
	`, s.Listener.Addr().String(),
		fixtureClientCertPEMBlock,
		fixtureClientKeyPEMBlock,
		fixtureCAPEMBlock,
		hcl,
	)
}

func testServerTLSConfig() (*tls.Config, error) {
	cert, err := tls.X509KeyPair(fixtureServerCertPEMBlock, fixtureServeKeyPEMBlock)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	block, _ := pem.Decode(fixtureCAPEMBlock)
	ca, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	pool.AddCert(ca)

	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pool,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
	}, nil
}

var (
	fixtureCAPEMBlock = []byte(`-----BEGIN CERTIFICATE-----
MIIFDTCCAvWgAwIBAgIJALLjNhWBWohjMA0GCSqGSIb3DQEBCwUAMBIxEDAOBgNV
BAMMB2Zha2UtY2EwHhcNMTcwNTIwMTkxMTEzWhcNMjcwNTE4MTkxMTEzWjASMRAw
DgYDVQQDDAdmYWtlLWNhMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA
rDNUVXlN/LRekUjJYQBK+7NBggbu6YUPOV7N7m1ncxGbduiLXnG1ZSEDERFgKWAg
d/gsHxEY9GQcICdo00MVFraqTLnMkpt3bLG9KRsTlURg+KTwbSYN7Sww3bd42QWT
if39HYIA4UYA1OMw7jGLAVQ5kScvfL4qA7AatyUgk2kOGGbtZm5cyhrlSM/HVUtL
OG8XSWVsY/xgVo0i4cDTLcxjWM0LKKHPnYymyVm/qasiLhRyS7SEL6FmzHiCjjGF
icnIi1vuXqXVEixaSWRGXUyZY7O0zb8cpXiYR9kWZtWzXpA9Qm38rWG/8Lsc/46/
8tKg0wSBYs93pzoOOjsaxZ2VVsTWM6GaGoca1ALzl3wnG0Mkt2Huz1gRe4a4TSdt
UBbpqG8hZyLhTsZXL8jIDRUt/V493TgFw6MccbYV6AgCquTYEDdeSQYJxC2Hk9XM
JqFRMZgTv4itoOVj7RAazYX4vp9vn7iyTtwRjotAP5dwXvsQ5qFsJCeKOid8Jqb9
Nk6eO0qEnYFeHV5ZUTppHivlrTzWfnNLRStGTYbzXY8GckUn7IWlfKU5wJIJoGh+
doW59bK537A+A2QnLuInQpg+1aTGndEjelbVQbmBRegnEHtcEIlA514twwExVESt
75bQYYiHHKhhGaF7BbDDh2yqzeBZMbFWY/Kii/PZ7X0CAwEAAaNmMGQwHQYDVR0O
BBYEFIlpEtDKJGJs5B4cmj82BejqGKaXMB8GA1UdIwQYMBaAFIlpEtDKJGJs5B4c
mj82BejqGKaXMBIGA1UdEwEB/wQIMAYBAf8CAQAwDgYDVR0PAQH/BAQDAgGGMA0G
CSqGSIb3DQEBCwUAA4ICAQAxxTSvSb2K6PNiGapGCqQ14qhT5y9gPQ0TYYTjkIj/
x1mTqGVNTlxGXlnY9rP1pfWuEettYvSBPKv86WPjUocwg80N8T7inFlpt9UROdh3
/thQOxzi35I+b9DMmWLcC6rKPZZ7PhGszM8no00OI2oMWML85EIjrENL1cgKKbeb
bZON36hHdQRmYMWJwY2rf9ayFardpulrL308c5y6J3Gq5hjs7rCv+Lc8PeBCCBEq
ECqnx8pwrv/W0VZlradAUl+Pdl4ti8Zg7/t/HH7JrUF68rQBYAHvjYi+7pt4S+96
iaQEkdankIwbPUO/gQ421cjXAemRpvtt3t9vazUOQwOYRWJLlxAeRHGOglI2NmxX
lN0mXfkNXhDZ5FDYxM9DRfi7eREcadxoRk+wROwW4WmoWYexMEWky/EdAlHWjXoo
qwb6lOJe6GoIbXJM6hGRFTVG8bb2B2zTDYwx51Cl5j1pDWmV8jSBWuiRVGRxUL1o
8XOLrhitBI1Qlnlfcx4eXv7uzseA4N0ncehxCxuHRc1KQussbpPKFU9U7AeusYQm
MI/UxunlbRzDzLzYunfqUWyvphXpjGJm8uN4jzrUCGC6Tv1XXO/7yjBt1APQNtP/
gstoe09eMCyhEfrxuHfEPqnrxpEgINk9fwlL8yfWkqPWxCZ3T6NVzLY0pSgSpy0r
tw==
-----END CERTIFICATE-----`)

	fixtureServerCertPEMBlock = []byte(`-----BEGIN CERTIFICATE-----
MIIErDCCApSgAwIBAgICEAAwDQYJKoZIhvcNAQELBQAwEjEQMA4GA1UEAwwHZmFr
ZS1jYTAeFw0xNzA1MjAxOTExMTNaFw0xODA1MjAxOTExMTNaMBYxFDASBgNVBAMM
C2Zha2Utc2VydmVyMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAp13f
kkHuHJRmmr53AxnlHd2kCc+g8YApdmb44IK621+1XDCtPM0ymljuXus/+iD1FRV/
IJ6vc6KxQKLTBFl2wJCwOojoPYRZiTW4kSWZ3tGlLWTKiJbNZ10+IDx4bt8pOvWg
j0bhj17d+IU6LUd09aGKb18GS0cbdKZYpOmih3o56J89sVp7MDV0ZvqJ4UFBNKIh
IRtl76Bg1A9wLnob0qPjVpEVozDmOl4XmziRHokmOv/JX7VgD9oqeWcY0BL1oYK/
NH7doflmGxnPt9dbni7wU+Vra3/5PzecUpcGurFW5Re23qpibK5V6cJQ2OPeOzFa
TadetPU33iZl4zRUCwIDAQABo4IBBjCCAQIwCQYDVR0TBAIwADARBglghkgBhvhC
AQEEBAMCBkAwMwYJYIZIAYb4QgENBCYWJE9wZW5TU0wgR2VuZXJhdGVkIFNlcnZl
ciBDZXJ0aWZpY2F0ZTAdBgNVHQ4EFgQUKknSKvnTFgNIguDzV/4vIReJzCEwQgYD
VR0jBDswOYAUiWkS0MokYmzkHhyaPzYF6OoYppehFqQUMBIxEDAOBgNVBAMMB2Zh
a2UtY2GCCQCy4zYVgVqIYzAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYB
BQUHAwEwJQYDVR0RBB4wHIIUbWF0Y2hib3guZXhhbXBsZS5jb22HBH8AAAEwDQYJ
KoZIhvcNAQELBQADggIBAKS9x7URH1ExonbDuuUpmLJE/C8GOyOfYtzt6WarDHpN
ltyqr7lN8o7tIvM/C3z5RwIEaN/OfsCHn/K5/fBs58qViULdv9aQtsa18xFcj8HT
jGSAFy80W/shbKTKcnrT3WVE+34HUh9JvXWa1TegY517p5a9y0TOX1vthT84Ynr/
r9ON5PJLRCrhXKs/itppHDctpqojqq8RB1BuTi9dSV8O5iCfr8Nu0zsUah/Y1p1u
6pPyeY0+qmv5Tdv6slubvgcMje5N+c9Y6x6BEws6anoVY+aWz+aLKOLDUqYRRlYQ
3nJxhcCNbzRmyK4OKi7mJRBHQOeIrkH2AWak1cGOtHMYly4OWCOboCPzIKPQGPtN
Ohl7Arp8bcjRdALNJP5qIz5sRgKvc4JJEAIAVVOlSzpycqBZJvuqa1PdY/cbnlRq
T57UdcLXSZQYV/wssqaJlfGctN5KkHh/jLRO9IZIg2ypTCna6yauoo1WGHXpffQi
zwWzHdVBak43lebsJTAQwiXr3RJsFZdfrx7rrivsFKc+zn/FDoonirX50tGXKgd/
cMuXpsptToBM8ij5xINDQRRHZLBYtWAhV3Ek7QxWFDFBCcEV+7lQp9ERwLuoEbNY
FQgkFiwjRRvIZ9ynCqiZXabS2bNl9Q1IKjxnCOLQ3oKWax9OcEeYUgR5dbf7IGHr
-----END CERTIFICATE-----`)

	fixtureServeKeyPEMBlock = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAp13fkkHuHJRmmr53AxnlHd2kCc+g8YApdmb44IK621+1XDCt
PM0ymljuXus/+iD1FRV/IJ6vc6KxQKLTBFl2wJCwOojoPYRZiTW4kSWZ3tGlLWTK
iJbNZ10+IDx4bt8pOvWgj0bhj17d+IU6LUd09aGKb18GS0cbdKZYpOmih3o56J89
sVp7MDV0ZvqJ4UFBNKIhIRtl76Bg1A9wLnob0qPjVpEVozDmOl4XmziRHokmOv/J
X7VgD9oqeWcY0BL1oYK/NH7doflmGxnPt9dbni7wU+Vra3/5PzecUpcGurFW5Re2
3qpibK5V6cJQ2OPeOzFaTadetPU33iZl4zRUCwIDAQABAoIBAFPC4GRTSLbW8m7Z
ibhsmkUTKsiaOAMFUDrol//MjXXC9YIY/mpii8PBZDLu64rkOaP+qSwLHuXxc2JU
2uTfXVZMU1ZINGqtNR49W4yQ0+w24cLRIaewSUZE3RXHDcL3Pqw6R8vM/pABO3fo
PVBx5bAU07KfTQgZoz0DD3QhVW0VSDDy7cUFmgg+I/qwhQwcUJz9Q0Fj1xclUpJU
4HtPoBX4JTrGw8QpVGeSr3KEzbAay6PGfzaaXxvAyLbV9IejrTthILEj/DwER9WV
3pMrbwiUWT3fge/jHrklewBsICWW1DcEOe5M5eUldvGAYBJz1sfaIkL1VKRbIEi2
1V8U9bkCgYEA2knjIibf48I4Ef+R5zcgSpHA1+2ScqV4UHR530VJjOOjnC36EdL5
mL11E+Bx1kwdHBMooMFLp88P5PtC7rJwq/XwxKG+70v1qKt3o/13Y6ua9TqEh8OI
gOLzaRW5ZE8ElNjE3Zl2A1GkLlHlmFUdx7rPnw/YoSykDver4+k2ue0CgYEAxEfk
5jl9dqA+3hnKb4FCIyd9auYRhUDVfu9/0I1mkw7QHG0Yb2JPunC3yZxxkwg+Cc7I
dDPoWCnQDSiLjlgrsqYImbOsmkRe+6LWXMRrhL5z1z2enRZEpHMORgCD/v3TaRQu
EEjT1hgiI5rlWUcu8KQrPDZjWlb3k1EGd+iZJtcCgYA04fORwYM6BUJaMeUh87vx
9M+YQCjbd3TnYOBpk7qW1Es9ufG8QbVQKI3li9loRjZDJ+0OzOVMOSCro6d6dmZP
cpyqtliwVmGkRC4O34f98IPw5wVWcqtuNg0sJyQrxezhNoay/MuXUD8LLbIGrpAx
Y/OKoGcl3M++BIhzBXvJnQKBgFuDmbmtvE1+0VEEfVoXzhpN4y/gLPMQE3qnd9Ro
2RZfpbBbPTVRhRLMUyRxCJMhGKvB+bwUJ5RTimlYKhkoCte0ifX/y83xasewWHnQ
KsEtex0z4awkIcT60ADbZK+S8OrhOcjl6766adBn+97wTXZtVKsyQIhyW+QXtwhZ
Lm7pAoGACNXFKjAJkEr32XyQJWJN3CUjg+92nnIdctbaNBzXHMFloDnDzRYT3PQr
N3SDfpFlvUSRxGyYJMHsEjzshAqGqRvCDunQrlPeChN/nFFbmK8nlDt9KtsB75/t
8qjhYFpzERomZfnK6TQwQsiXrPgpLrlK6BCkWdq47h7rlwISMtA=
-----END RSA PRIVATE KEY-----`)

	fixtureClientCertPEMBlock = []byte(`-----BEGIN CERTIFICATE-----
MIIEYDCCAkigAwIBAgICEAEwDQYJKoZIhvcNAQELBQAwEjEQMA4GA1UEAwwHZmFr
ZS1jYTAeFw0xNzA1MjAxOTExMTRaFw0xODA1MjAxOTExMTRaMBYxFDASBgNVBAMM
C2Zha2UtY2xpZW50MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAurFI
FX9Il6sPgQTrEU4CHvcpfTqEVk9XPlEf6oP3qsjDjYnf+KT5IgHzuAxEEqVyStCb
P8eY+1DYUDPfe8bT7FNCO9zDuFSTA5XybeXUp1q1Yo8n+wRneFyFc0X3EllGZkV/
QGnmNpOkRVhwdaf6pqxK6sMxG1i5EPV0+4soAxG98ci2rmI6EKWtQBdUy1X8GizF
EiQPfpNRljkNU2RK5EXf6eCEGACjsQMnccXz5u0bXNOnw65zjLoA69RTNaRkt/TM
DTN6svXFpiBNeEeEOk9M4UTJRw2K6at4sPQ/2GMd061Vmy0iMibfc1MleD1oyvqp
XR+gBPtP49Yc2GPpkwIDAQABo4G7MIG4MAkGA1UdEwQCMAAwEQYJYIZIAYb4QgEB
BAQDAgeAMDMGCWCGSAGG+EIBDQQmFiRPcGVuU1NMIEdlbmVyYXRlZCBDbGllbnQg
Q2VydGlmaWNhdGUwHQYDVR0OBBYEFErCDJadACeDhkM24qOKnOJIlglKMB8GA1Ud
IwQYMBaAFIlpEtDKJGJs5B4cmj82BejqGKaXMA4GA1UdDwEB/wQEAwIF4DATBgNV
HSUEDDAKBggrBgEFBQcDAjANBgkqhkiG9w0BAQsFAAOCAgEAgDFI4cG9hyY+BvDS
QD9wesNcdfrdR8cA1fQm5z866U3xFZwG6+xX449TJv9um3VFPXJL+0k9kqLBeN/V
ez1YFlKdv59Vylnzd9b+WG5CJjd8awNWs2Fe7fxTlUkpW5mspwfaua9sFPD+Dc9K
j/Te2SIENLo1PfBKjelot1T3CIOs3HbdtfLbAsS29Swqog8FJaNKQd1dAkLD+Smv
tkTHgVpdqbJyAnKBqEckFgfSvN5Hfy/bsi/ugmWLfANdSdfCyIsi0g+h3CJiyWot
VQkfhbuc4b7j0O51rfrHMF9mp57tp0pVZ8bLS5N+uVWYeYJUAbXVFIK868RCSTqb
dWW0SnrWptNyWzybsX8jugy+sFAazpjvmzuM74pZx2mIVnNLO5o46kW+rry6/e76
GcLy80dPGVJG6HVTNYR+Y0g35CIdzHgs8p7EEd2N3MfoTagPXaz0JEx0sOWxuHaJ
88rFzFw/NDvhgcc1dQ8fAK74JBwsZ/ufrk0A637UDWbacPLNNDyeaGqt+ab6x05W
lpE2FvmAsFAq0/bi6H9VvTncbi8FbthhqDKmGS05J7z1oTW9a2U9h/l5H2SFFpZI
IRI7M4wmDZAyT3SN4LyyrZUld8UjRgHmHATWfSTEYDP5vv9rvuG391mPo99VzhzR
fkdKG26l+ZmAvPPMWumumYvBu/0=
-----END CERTIFICATE-----`)

	fixtureClientKeyPEMBlock = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAurFIFX9Il6sPgQTrEU4CHvcpfTqEVk9XPlEf6oP3qsjDjYnf
+KT5IgHzuAxEEqVyStCbP8eY+1DYUDPfe8bT7FNCO9zDuFSTA5XybeXUp1q1Yo8n
+wRneFyFc0X3EllGZkV/QGnmNpOkRVhwdaf6pqxK6sMxG1i5EPV0+4soAxG98ci2
rmI6EKWtQBdUy1X8GizFEiQPfpNRljkNU2RK5EXf6eCEGACjsQMnccXz5u0bXNOn
w65zjLoA69RTNaRkt/TMDTN6svXFpiBNeEeEOk9M4UTJRw2K6at4sPQ/2GMd061V
my0iMibfc1MleD1oyvqpXR+gBPtP49Yc2GPpkwIDAQABAoIBACeSGghMcVOMc33S
UAzb7wEnPEkJ1TECIijYQx6PGDi/0ws2FR37wb6ekU0KdIdLQB1xd+ad5OQn76GY
TR9MNnEZ+Kj9kxKIAp049CitFVTfmiCo3T2MYm4VlkenpcXi3FQjGOTLTXt18dSs
+TFHCI65aCu4cbktJhTdIg2LIlD73cZbhw5c3XqzNrauzCj+bdgwomqKmUkiTNMt
OVxmLRiuRoe13imChNVm6m1TQ+jUeEVdbEXnThTKrq+DtQdwcS6l9ucl3P2ZlhMT
HoBI72bvUJ6CC+pc4Y/AKrUW2tFMKk02W2cSNWqLwXwmS5WkZwl4Jwui+7DMsbqB
X3Rp3okCgYEA7CAxeGc2YnWru3ck789n4lkT7s6G8iH5CuJd/UZ4fHeMn2rfO0dB
FhiOaLXOcfpU5YypiZyJjy+584qQktyBCeWJXTb1USpKOsamHddt65rYCqOnKAaH
W3WqmhZDVneNQWW4UsiQut+0ZYBfmZQamSBWwBdn6cUtvRxxfoj2jAUCgYEAymfy
9TCY0SUNthz0xlS6woR5fKh5+TCs0RQhApE25W1OeCJV7p2v06DA+KvUGcdvL1j4
WOF4mQXfDiHE0ceOo1uMcqsH2U2TUGdCkxVQgdQRaKaAogl2IoqCrSdsq4btyh9T
cwe7Jyf29psaeUwYZ5SwBzSKhHnZED2+VR0JKrcCgYBn2AipoQqj5oguG8ncxWQ0
gWRow99JIXO7O66GMrXOV206tu+RzFZtd0M5/arbKXKouWHeKT+9/wlSd//49oyx
Y4czvXXJykV27+IigZnP4ftdQnfC/IwOxwLOXTgkENPIjQmxLo+n/7YAZaKlkiLY
cQZ12FVU0+i3oIixU17KWQKBgDhMg5rJoqgB63dtRHRqGuyCFpyi7BJxBJC/TZM+
OwvDxKDLxCUz/TUbMLG6caud+oIr+CAYzweZR3rRz8IeBMHRdBZtFijOWBx0LGNm
+VazWwhFz9/CS/a9mi15mtN3G2suHXMQgnEYv6vGZq24ic094VyPs7u3fLX0xp08
D1GvAoGAJsnXun/4Rxr9VheiAry3D8heqUu3BykDexn1w66Y1Nj/+uwnOvEYa5ko
hL+wXhSF6ibNgPKjpw4sXJgkC6RF8OVNslz9n8H4tZ14K9ba9mTfvSDNab7KVT6m
ldAiovW49XtR31J7RwIE0DjKt0hoWAkq/L0n9rQ0vdqIl3to/E4=
-----END RSA PRIVATE KEY-----`)
)
