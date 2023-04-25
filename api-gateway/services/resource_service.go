package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ResourceService struct {
	baseURL string
}

func NewResourceService(resourceManagerURL string) *ResourceService {
	return &ResourceService{baseURL: resourceManagerURL}
}

func (rs *ResourceService) ProxyRequest(w http.ResponseWriter, r *http.Request, endpoint string, method string) error {
	url := fmt.Sprintf("%s%s", rs.baseURL, endpoint)
	fmt.Println("URL:", url)
	fmt.Println("Method:", method)

	client := &http.Client{}
	newRequest, err := http.NewRequest(method, url, r.Body)
	if err != nil {
		return err
	}

	// Copy headers from the original request to the new request
	for k, v := range r.Header {
		newRequest.Header.Set(k, v[0])
	}

	resp, err := client.Do(newRequest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Response status:", resp.Status)
	fmt.Println("Response body:", string(body))
	fmt.Println("Response headers:", resp.Header)

	// Copy the status code from the resource-manager response
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
