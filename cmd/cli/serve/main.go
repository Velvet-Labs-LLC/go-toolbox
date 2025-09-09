//nolint:cyclop // This package contains complex server setup logic
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"
)

//nolint:cyclop,gocognit,nestif // This function handles CLI argument parsing and server setup
func main() {
	tlsFlag := flag.Bool("tls", false, "Enable HTTPS (default: false)")
	certFile := flag.String("cert", "", "Path to TLS certificate file (optional)")
	keyFile := flag.String("key", "", "Path to TLS key file (optional)")
	port := flag.Int("port", 0, "Port to listen on (default: 8080 for HTTP, 8443 for HTTPS)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--tls] [--cert CERT] [--key KEY] [--port PORT] <directory>\n", os.Args[0])
		os.Exit(1)
	}
	dir := flag.Arg(0)

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatalf("Directory does not exist: %s", dir)
	}

	// Get local IP address for user info
	ip := getLocalIP()
	if ip == "" {
		ip = "localhost"
	}

	var listenPort int
	var scheme string
	if *tlsFlag {
		scheme = "https"
		if *port != 0 {
			listenPort = *port
		} else {
			listenPort = 8443
		}
	} else {
		scheme = "http"
		if *port != 0 {
			listenPort = *port
		} else {
			listenPort = 8080
		}
	}

	fmt.Printf("Serving %s on %s://%s:%d (accessible on your LAN)\n", dir, scheme, ip, listenPort)
	http.Handle("/", http.FileServer(http.Dir(dir)))

	addr := fmt.Sprintf("0.0.0.0:%d", listenPort)
	if *tlsFlag {
		var cert, key string
		if *certFile != "" && *keyFile != "" {
			cert = *certFile
			key = *keyFile
		} else {
			// Generate self-signed cert
			cert, key = generateSelfSignedCert()
		}
		srv := &http.Server{
			Addr:              addr,
			Handler:           nil,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      15 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			IdleTimeout:       60 * time.Second,
		}
		if cert == "" || key == "" {
			log.Fatal("Failed to generate or find certificate and key for HTTPS.")
		}
		if *certFile == "" && *keyFile == "" {
			// Use in-memory cert
			tlsCert, err := tls.X509KeyPair([]byte(cert), []byte(key))
			if err != nil {
				log.Fatalf("Failed to load self-signed certificate: %v", err)
			}
			srv.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{tlsCert},
				MinVersion:   tls.VersionTLS12,
			}
			log.Fatal(srv.ListenAndServeTLS("", ""))
		}
		log.Fatal(srv.ListenAndServeTLS(cert, key))
	}
	srv := &http.Server{
		Addr:              addr,
		Handler:           nil,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
} // generateSelfSignedCert returns PEM-encoded cert and key as strings
func generateSelfSignedCert() (string, string) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Printf("Failed to generate private key: %v", err)
		return "", ""
	}
	tmpl := x509.Certificate{
		SerialNumber:          big.NewInt(time.Now().UnixNano()),
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
	if err != nil {
		log.Printf("Failed to create certificate: %v", err)
		return "", ""
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	return string(certPEM), string(keyPEM)
}

// getLocalIP returns the first non-loopback, non-link-local IPv4 address
func getLocalIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, addrErr := iface.Addrs()
		if addrErr != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() || ip.IsLinkLocalUnicast() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			return ip.String()
		}
	}
	return ""
}
