package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const (
	itemsDir = "/items"
	certsDir  = "/certs"
	port      = ":41596"
)

type Metadata struct {
	Name        string `json:"name"`
	Ext         string `json:"ext"`
	NoThumbnail bool   `json:"noThumbnail"`
}

func main() {
	if _, err := os.Stat(itemsDir); os.IsNotExist(err) {
		log.Fatalf("STORAGE_DIR: `%s` not found", itemsDir)
	}

	http.HandleFunc("/api/item", serveItem)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	startServer()
}

func serveItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}

	thumbnail := false
	if thumbnailValue := r.URL.Query().Get("thumbnail"); thumbnailValue != "" {
		var err error
		thumbnail, err = strconv.ParseBool(thumbnailValue)
		if err != nil {
			http.Error(w, "Invalid thumbnail value", http.StatusBadRequest)
			return
		}
	}

	itemDir := filepath.Join(itemsDir, id+".info")
	metadataPath := filepath.Join(itemDir, "metadata.json")

	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	metadataBytes, err := os.ReadFile(metadataPath)
	if err != nil {
		log.Printf("Error reading metadata: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var metadata Metadata
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		log.Printf("Error parsing metadata: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if metadata.Name == "" {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	originalPath := filepath.Join(itemDir, metadata.Name+"."+metadata.Ext)
	itemPath := originalPath
	if thumbnail && !metadata.NoThumbnail {
		thumbnailPath := filepath.Join(itemDir, metadata.Name+"_thumbnail.png")
		if fileExists(thumbnailPath) {
			itemPath = thumbnailPath
		}
	}

	if _, err := os.Stat(itemPath); os.IsNotExist(err) {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, itemPath)
}

func startServer() {
	certPath := filepath.Join(certsDir, "cert.pem")
	keyPath := filepath.Join(certsDir, "key.pem")

	certExists := fileExists(certPath)
	keyExists := fileExists(keyPath)

	if certExists && keyExists {
		log.Printf("Starting HTTPS server on https://0.0.0.0%s", port)
		if err := http.ListenAndServeTLS(port, certPath, keyPath, nil); err != nil {
			log.Fatalf("HTTPS server failed: %v", err)
		}
	} else {
		log.Printf("Starting HTTP server on http://0.0.0.0%s", port)
		if err := http.ListenAndServe(port, nil); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}