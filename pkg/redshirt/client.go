package redshirt

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/pantheon-systems/certinel"
	"github.com/pantheon-systems/certinel/pollwatcher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var interval = 10 * time.Second

// NewTLSConnection will initiate a connection to riker with given address caFile and client .pem
func NewTLSConnection(addr, caFile, tlsFile, duration string) (*grpc.ClientConn, error) {
	watcher := pollwatcher.New(tlsFile, tlsFile, interval)

	sentinel := certinel.New(watcher, nil, func(err error) {
		fmt.Printf("certinel was unable to reload the certificate: %s", err)
	})
	sentinel.Watch()

	tlsConfig := &tls.Config{
		GetCertificate: sentinel.GetCertificate,
	}

	return grpc.Dial(addr, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
}
