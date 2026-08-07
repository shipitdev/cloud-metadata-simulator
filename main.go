package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Metadata struct {
	InstanceID       string   `yaml:"instance_id" json:"instance_id"`
	Hostname         string   `yaml:"hostname" json:"hostname"`
	LocalIPv4        string   `yaml:"local_ipv4" json:"local_ipv4"`
	PublicIPv4       string   `yaml:"public_ipv4" json:"public_ipv4"`
	AMIID            string   `yaml:"ami_id" json:"ami_id"`
	AvailabilityZone string   `yaml:"availability_zone" json:"availability_zone"`
	PublicKeys       []string `yaml:"public_keys" json:"public_keys"`
	UserData         string   `yaml:"user_data" json:"user_data"`
}

var globalMeta Metadata

func loadMetadata(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read metadata file: %w", err)
	}
	var meta Metadata
	if err := yaml.Unmarshal(data, &meta); err != nil {
		return fmt.Errorf("failed to parse metadata YAML: %w", err)
	}
	globalMeta = meta
	return nil
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// IMDS Root & Leaf Routes
	mux.HandleFunc("/latest/meta-data/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/latest/meta-data/")

		switch path {
		case "", "/":
			fmt.Fprintln(w, "ami-id\navailability-zone\nhostname\ninstance-id\nlocal-ipv4\npublic-ipv4\npublic-keys/")
		case "instance-id":
			fmt.Fprint(w, globalMeta.InstanceID)
		case "hostname":
			fmt.Fprint(w, globalMeta.Hostname)
		case "local-ipv4":
			fmt.Fprint(w, globalMeta.LocalIPv4)
		case "public-ipv4":
			fmt.Fprint(w, globalMeta.PublicIPv4)
		case "ami-id":
			fmt.Fprint(w, globalMeta.AMIID)
		case "availability-zone":
			fmt.Fprint(w, globalMeta.AvailabilityZone)
		case "public-keys", "public-keys/":
			fmt.Fprintln(w, "0=key-0")
		case "public-keys/0/openssh-key", "public-keys/0":
			if len(globalMeta.PublicKeys) > 0 {
				fmt.Fprint(w, globalMeta.PublicKeys[0])
			} else {
				http.Error(w, "404 page not found", http.StatusNotFound)
			}
		default:
			http.Error(w, "404 page not found", http.StatusNotFound)
		}
	})

	// Raw User Data Route (#cloud-config)
	mux.HandleFunc("/latest/user-data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, globalMeta.UserData)
	})

	// Structured JSON Route
	mux.HandleFunc("/latest/meta-data.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(globalMeta)
	})

	// Health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	return mux
}

func main() {
	configPath := flag.String("config", "metadata.yaml", "Path to metadata YAML configuration file")
	port := flag.Int("port", 8080, "HTTP server listening port")
	flag.Parse()

	if err := loadMetadata(*configPath); err != nil {
		log.Fatalf("Error loading metadata: %v", err)
	}

	mux := setupRoutes()
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting cloud metadata simulator on http://localhost%s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
