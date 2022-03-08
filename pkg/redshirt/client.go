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

// NewTLSConnection will initiate a connection to riker with given address caFile and client .pem
func NewTLSConnection(addr, caFile, tlsFile, duration string) (*grpc.ClientConn, error) {
	interval, err := time.ParseDuration(duration)
	if err != nil {
		return nil, fmt.Errorf("Could not parse duration (%s) for tls: %s", duration, err.Error())
	}

	watcher := pollwatcher.New(tlsFile, caFile, interval)

	sentinel := certinel.New(watcher, nil, func(err error) {
		fmt.Errorf("certinel was unable to reload the certificate: %s", err)
	})
	sentinel.Watch()

	tlsConfig := &tls.Config{
		GetCertificate: sentinel.GetCertificate,
	}

	return grpc.Dial(addr, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
}
