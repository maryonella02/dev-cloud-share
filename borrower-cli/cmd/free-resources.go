package cmd

import (
	"borrower-cli/helpers"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(freeResourcesCmd)
}

var freeResourcesCmd = &cobra.Command{
	Use:   "free-resources",
	Short: "Get all free available resources",
	Run: func(cmd *cobra.Command, args []string) {
		// Make a GET request to the API endpoint
		apiURL := "https://localhost:8440/api/v1/resources/free"

		// Create a custom HTTP client with insecure TLS configuration
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			log.Fatalf(fmt.Sprintf("%s", err))
		}
		req.Header.Set("Content-Type", "application/json")
		var token string
		if Token == "" {
			token, err = helpers.GetToken()
			if err != nil {
				fmt.Println("Error reading token:", err)
			}
		} else {
			token = Token
		}

		req.Header.Add("Authorization", "Bearer "+token)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf(fmt.Sprintf("%s", err))
		}
		defer resp.Body.Close()

		// Check the response status code
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Error getting free resources: status code %d", resp.StatusCode)
		}

		// Decode the response body into a slice of Resource structs
		var resources []Resource
		err = json.NewDecoder(resp.Body).Decode(&resources)
		if err != nil {
			log.Fatalf("Error decoding response: %v", err)
		}

		// Print the list of free resources
		if len(resources) > 0 {
			fmt.Println("Free Resources:")
			for i, resource := range resources {
				fmt.Print("\n" + strconv.Itoa(i+1) + ".")
				fmt.Print("\tCPU Cores: " + strconv.Itoa(resource.CPUCores))
				fmt.Print("\tMemory(MB): " + strconv.Itoa(resource.MemoryMB) + "\n")
			}
		} else {
			fmt.Println("No resources available.")
		}

	},
}
