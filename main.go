package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net"
	"os"
	"time"
)

// RawAPIServer functions as the TLS encrypted endpoint for beacon data
type RawAPIServer struct {
	*tls.Conn
	cert []byte
	key  []byte
	ca   []byte
	addr string
	// CAKey []byte // Beacon is assumed not to have the CAKey
}

// setupTLSConfig returns a tls.Config for a credential set
func setupTLSConfig(cert []byte, key []byte, ca []byte) (*tls.Config, error) {
	// TLS config
	var tlsConfig tls.Config

	//Use only modern ciphers
	tlsConfig.CipherSuites = []uint16{
		tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	}

	// Use only TLS v1.2
	tlsConfig.MinVersion = tls.VersionTLS12

	// Don't allow session resumption
	tlsConfig.SessionTicketsDisabled = true

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(ca)
	tlsConfig.RootCAs = certPool
	tlsConfig.ClientCAs = certPool

	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert

	keypair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return &tlsConfig, err
	}
	tlsConfig.Certificates = []tls.Certificate{keypair}

	return &tlsConfig, nil
}

const rawNetwork = "tcp4"

// NewRawAPIServer instantiates the encrypted server we'll send data to
func NewRawAPIServer(cert []byte, key []byte, ca []byte, addr string) (listener net.Listener, err error) {
	var tlsConfig *tls.Config

	tlsConfig, err = setupTLSConfig(cert, key, ca)
	if err != nil {
		return nil, err
	}
	listener, err = tls.Listen(rawNetwork, addr, tlsConfig)
	return listener, err
}

func handleAPIData(conn net.Conn) {
	var err error

	t := time.Now()
	log.Printf("%v\n", t.UTC())

	for {
		_, err = io.Copy(conn, conn)
		if err != nil {
			// TODO: Write connection errors to debug
			log.Printf("Connection error: %v\n", err)
			conn.Close()
			return
		}
	}
}

// RawAPIAddr is the address the API
const RawAPIAddr = "0.0.0.0:27001"

func main() {
	var rawAPIServer net.Listener
	var err error

	cert := []byte(os.Getenv("CERT"))
	key := []byte(os.Getenv("KEY"))
	ca := []byte(os.Getenv("CA"))

	rawAPIServer, err = NewRawAPIServer(cert, key, ca, RawAPIAddr)
	if err != nil {
		log.Fatalf("Unable to create the raw API server: %v\n", err)
	}

	log.Println("Beacon up!")
	for {
		var conn net.Conn
		conn, err = rawAPIServer.Accept()
		if err != nil {
			log.Printf("Error with request: %v\n", err)
		}
		go handleAPIData(conn)
	}
}
