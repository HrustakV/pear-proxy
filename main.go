package main

import (
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	Address    string   `yaml:"address"`
	Endpoint   string   `yaml:"endpoint"`
	AllowedIPs []string `yaml:"allowed_ips"`
}

var config Config

var prefix = "[PEAR-PROXY] "

func handleProxy(w http.ResponseWriter, r *http.Request) {
	targetUrl, err := url.Parse(config.Endpoint)
	if err != nil {
		http.Error(w, prefix+"Invalid target URL", http.StatusInternalServerError)
		return
	}

	allowedIPs := config.AllowedIPs

	requestIP := strings.Split(r.RemoteAddr, ":")[0]

	authorizedIP := slices.Contains(allowedIPs, requestIP)

	if !authorizedIP {
		http.Error(w, prefix+"IP address not authorized", http.StatusUnauthorized)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	proxy.ServeHTTP(w, r)
}

func main() {
	configFile, err := os.ReadFile("/etc/pear-proxy/config.yaml")
	if err != nil {
		log.Fatalf("%vError reading config %v", prefix, err)
	}

	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("%vError parsing YAML: %v", prefix, err)
	}

	http.HandleFunc("/", handleProxy)

	log.Printf("%vProxy listening on %v", prefix, config.Address)

	if err := http.ListenAndServe(config.Address, nil); err != nil {
		log.Fatal(err)
	}

}
