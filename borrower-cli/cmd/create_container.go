package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
)

var (
	image    string
	cpuCores int
	memoryMB int
)

var createContainerCmd = &cobra.Command{
	Use:   "create-container",
	Short: "Create a new container",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating container with image:", image)
		createContainer(image, cpuCores, memoryMB)
	},
}

func init() {
	createContainerCmd.Flags().StringVarP(&image, "image", "i", "", "Image to use for creating the container (required)")
	createContainerCmd.Flags().IntVar(&cpuCores, "cpu-cores", 0, "CPU cores for the container")
	createContainerCmd.Flags().IntVar(&memoryMB, "memory-mb", 0, "Memory limit in MB for the container")
	createContainerCmd.MarkFlagRequired("image")
	rootCmd.AddCommand(createContainerCmd)
}

func createContainer(image string, cpuCores, memoryMB int) {
	// Convert CPU cores to nano-cpus (10^-9)
	nanoCPUs := int64(cpuCores * 1e9)

	// Convert memory limit to bytes
	memoryBytes := int64(memoryMB * 1024 * 1024)

	// Prepare the POST request
	url := "https://localhost:8440/api/v1/containers"
	payload := map[string]interface{}{
		"image":     image,
		"nano-cpus": nanoCPUs,
		"memory":    memoryBytes,
	}
	// TODO: use defined values earlier when requesting resource, if no resource constrains specified

	data, _ := json.Marshal(payload)
	body := bytes.NewReader(data)
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", "application/json")

	// Create a custom HTTP client with insecure TLS configuration
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Read and print the response
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Response:", string(respBody))

	// Save container ID to a local file
	var result map[string]string
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return
	}
	containerID = result["id"]
	err = os.WriteFile("container_id.txt", []byte(containerID), 0644)
	if err != nil {
		fmt.Println("Error saving container ID:", err)
	} else {
		fmt.Println("Container ID saved to container_id.txt")
	}
}

// ./borrower-cli create-container --image nginx --cpu-cores 2 --memory-mb 1024

// The code will convert:
// cpu-cores value of 2 to 2000000000 nano-cpus
// memory-mb value of 1024 to 1073741824 bytes, and send them to the gateway
