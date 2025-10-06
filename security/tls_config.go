package security

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

func NewtlsConfig() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair("./tls/client.cert.pem", "./tls/client.key.pem")
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile("./tls/ca.cert.pem")
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}, nil
}
