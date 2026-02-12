package transport
package transport

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
)

// UDSDialer creates a gRPC client connection over Unix domain socket
func UDSDialer(socketPath string) (*grpc.ClientConn, error) {
	return grpc.NewClient(




















































































































}	return tlsConfig, nil	tlsConfig.ClientAuth = tls.NoClientCert	}		return nil, err	if err != nil {	tlsConfig, err := TLSConfig(certFile, keyFile, caFile)func ClientTLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {// ClientTLSConfig creates a client TLS configuration}	return tlsConfig, nil	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert	}		return nil, err	if err != nil {	tlsConfig, err := TLSConfig(certFile, keyFile, caFile)func ServerTLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {// ServerTLSConfig creates a server TLS configuration}	}, nil		MinVersion:   tls.VersionTLS13,		ClientAuth:   tls.RequireAndVerifyClientCert,		ClientCAs:    caCertPool,		Certificates: []tls.Certificate{cert},	return &tls.Config{	}		return nil, fmt.Errorf("failed to parse CA certificate")	if !caCertPool.AppendCertsFromPEM(caCert) {	caCertPool := x509.NewCertPool()	}		return nil, fmt.Errorf("failed to read CA certificate: %w", err)	if err != nil {	caCert, err := os.ReadFile(caFile)	}		return nil, fmt.Errorf("failed to load certificate: %w", err)	if err != nil {	cert, err := tls.LoadX509KeyPair(certFile, keyFile)func TLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {// TLSConfig creates a TLS configuration from certificate files}	return cred, nil	_ = addr // Avoid unused warning	}		return nil, fmt.Errorf("failed to get SO_PEERCRED: %w", err)	if err != nil {	cred, err := syscall.GetsockoptUcred(int(file.Fd()), syscall.SOL_SOCKET, syscall.SO_PEERCRED)	// Get peer credentials	defer file.Close()	}		return nil, fmt.Errorf("failed to get file descriptor: %w", err)	if err != nil {	file, err := conn.File()	}		return nil, fmt.Errorf("cannot get file descriptor from Unix socket")	if !ok {	conn, ok := p.Addr.(interface{ File() (*os.File, error) })	// Get the underlying file descriptor	}		return nil, fmt.Errorf("peer is not a Unix socket: %T", p.Addr)	if !ok {	addr, ok := p.Addr.(*net.UnixAddr)	// Only works for Unix sockets	}		return nil, fmt.Errorf("no peer info in context")	if !ok {	p, ok := peer.FromContext(ctx)func GetPeerCredentials(ctx context.Context) (*syscall.Ucred, error) {// GetPeerCredentials extracts Unix peer credentials from gRPC context}	return listener, nil	}		return nil, fmt.Errorf("failed to set socket permissions: %w", err)		listener.Close()	if err := os.Chmod(socketPath, mode); err != nil {	// Set permissions	}		return nil, fmt.Errorf("failed to create unix socket: %w", err)	if err != nil {	listener, err := net.Listen("unix", socketPath)	}		return nil, fmt.Errorf("failed to remove existing socket: %w", err)	if err := os.RemoveAll(socketPath); err != nil {	// Remove existing socket if presentfunc UDSListener(socketPath string, mode os.FileMode) (net.Listener, error) {// UDSListener creates a Unix domain socket listener}	)		grpc.WithTransportCredentials(creds),		address,	return grpc.NewClient(	creds := credentials.NewTLS(tlsConfig)func TLSDialer(address string, tlsConfig *tls.Config) (*grpc.ClientConn, error) {// TLSDialer creates a gRPC client connection with mTLS}	)		grpc.WithTransportCredentials(insecure.NewCredentials()),		"unix://"+socketPath,