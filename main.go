package main

import (
	"crypto/tls"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"HTTPs-Golang/enums"
	"HTTPs-Golang/logger"
)

func main() {
	log := logger.NewLogger("")

	server, err := enums.NewServer()
	if err != nil {
		log.Fatal("Error initializing server", map[string]interface{}{
			"error": err,
		})
	}

	port := "443"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

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

	log.Info("Server started", map[string]interface{}{
		"port": port,
	})

	log.Info("Login URL", map[string]interface{}{
		"url": server.Config.LoginURL,
	})

	certFile := filepath.Join("config", "SSL", "server.crt")
	keyFile := filepath.Join("config", "SSL", "server.key")

	if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil {
		log.Fatal("Error starting server", map[string]interface{}{
			"error": err,
		})
	}
}
