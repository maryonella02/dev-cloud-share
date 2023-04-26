package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
)

var image string

var createContainerCmd = &cobra.Command{
	Use:   "create-container",
	Short: "Create a new container",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating container with image:", image)
		createContainer(image)
	},
}

func init() {
	createContainerCmd.Flags().StringVarP(&image, "image", "i", "", "Image to use for creating the container (required)")
	createContainerCmd.MarkFlagRequired("image")
	rootCmd.AddCommand(createContainerCmd)
}
func createContainer(image string) {
	// Prepare the POST request
	url := "http://localhost:8081/api/v1/containers"
	payload := map[string]interface{}{
		"image": image,
	}
	data, _ := json.Marshal(payload)
	body := bytes.NewReader(data)
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", "application/json")

	// Send the POST request
	client := &http.Client{}
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

//./borrower-cli create-container --image nginx
