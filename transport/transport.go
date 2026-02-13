package transport

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
)

// UDSDialer creates a gRPC client connection over Unix domain socket
func UDSDialer(socketPath string) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		"unix:"+socketPath,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}

// UDSListener creates a Unix domain socket listener
func UDSListener(socketPath string, mode os.FileMode) (net.Listener, error) {
	// Remove existing socket if present
	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to remove existing socket: %w", err)
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on unix socket: %w", err)
	}

	if err := os.Chmod(socketPath, mode); err != nil {
		listener.Close()
		return nil, fmt.Errorf("failed to chmod socket: %w", err)
	}

	return listener, nil
}

// TLSDialer creates a gRPC client connection with TLS
func TLSDialer(address string, tlsConfig *tls.Config) (*grpc.ClientConn, error) {
	creds := credentials.NewTLS(tlsConfig)
	return grpc.NewClient(address, grpc.WithTransportCredentials(creds))
}

// TLSConfig creates a TLS configuration for mutual TLS
func TLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load cert/key pair: %w", err)
	}

	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA cert")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
		ClientCAs:    caPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
	}, nil
}

// ServerTLSConfig creates a server-side TLS configuration with mTLS
func ServerTLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
	return TLSConfig(certFile, keyFile, caFile)
}

// ClientTLSConfig creates a client-side TLS configuration for connecting to mTLS servers
func ClientTLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load cert/key pair: %w", err)
	}

	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA cert")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
		MinVersion:   tls.VersionTLS13,
	}, nil
}

// PeerInfo holds peer identity information extracted from gRPC context
type PeerInfo struct {
	Address   string
	IsUnix    bool
	IsTLS     bool
	TLSSubject string
}

// GetPeerInfo extracts peer information from a gRPC context
func GetPeerInfo(ctx context.Context) (*PeerInfo, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no peer info in context")
	}

	info := &PeerInfo{}

	if p.Addr != nil {
		info.Address = p.Addr.String()
		if _, ok := p.Addr.(*net.UnixAddr); ok {
			info.IsUnix = true
		}
	}

	if tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo); ok {
		info.IsTLS = true
		if len(tlsInfo.State.PeerCertificates) > 0 {
			info.TLSSubject = tlsInfo.State.PeerCertificates[0].Subject.CommonName
		}
	}

	return info, nil
}
