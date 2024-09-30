package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// The target URL to which we want to proxy the requests
	target := "https://mvp-ds.dev.mat.ptr.smartlivingnext.de"

	// Parse the URL
	targetURL, err := url.Parse(target)
	if err != nil {
		http.Error(w, "Error parsing target URL", http.StatusInternalServerError)
		return
	}

	// Create a new request based on the original request but with the target URL
	slog.Info(
		"Method", r.Method,
		"NewRequest", "targetURL",
		targetURL.String(),
		"RequestUri", r.RequestURI,
	)

	proxyReq, err := http.NewRequest(r.Method, targetURL.String()+r.RequestURI, r.Body)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Copy original headers to the new request
	proxyReq.Header = r.Header

	slog.Debug("proxyReq", "Header", proxyReq.Header)

	// Find X-Access-Token, which is set by oauth2-proxy
	token := proxyReq.Header.Get("X-Access-Token")
	// Set token as Bearer
	if token != "" {
		slog.Info("got", "token", token)
		proxyReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy response headers from the target response
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Set(key, value)
		}
	}

	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Copy the status code from the response
	w.WriteHeader(resp.StatusCode)

	// Copy the response body from the target response to the client
	io.Copy(w, resp.Body)
}

func main() {
	http.HandleFunc("/", proxyHandler)
	fmt.Println("starting proxy server on port 3333")
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		log.Fatal(err)
	}
}
