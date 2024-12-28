package main

import (
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"HTTPs-Golang/enums"
	"HTTPs-Golang/logger"
)

const (
	defaultPort = "443"
)

func main() {
	log := logger.NewLogger("")

	server, err := enums.NewServer()
	if err != nil {
		log.Fatal("Error initializing server", map[string]interface{}{
			"error": err,
		})
	}

	log.Debug("Server configuration loaded", map[string]interface{}{
		"config": server.Config,
	})

	port := getPort()

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
	}

	srv := &http.Server{
		Addr:         ":" + port,
		TLSConfig:    tlsConfig,
		Handler:      server,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := srv.ListenAndServeTLS("config/SSL/server.crt", "config/SSL/server.key"); err != nil {
		log.Fatal("Server failed", map[string]interface{}{
			"error": err,
		})
	}
}

func getPort() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return defaultPort
}
