package services

import (
	"crypto/tls"
	"io"
	"net/http"
)

type BaseService struct {
	ServiceURL string
}

func (bs *BaseService) ProxyRequest(w http.ResponseWriter, r *http.Request, endpoint, method string) error {
	target := bs.ServiceURL + endpoint
	req, err := http.NewRequest(method, target, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	// Copy headers from original request to the new request
	for k, v := range r.Header {
		req.Header[k] = v
	}

	// Create a custom HTTP client with insecure TLS configuration
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	defer resp.Body.Close()

	// Copy response headers and status code
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}
