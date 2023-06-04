package cmd

import (
	"borrower-cli/helpers"
	"crypto/tls"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
)

// Add a new stopContainerCmd
var stopContainerCmd = &cobra.Command{
	Use:   "stop-container",
	Short: "Stop a container",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stopping container...")
		stopContainer(containerID)
	},
}

// Modify the init function to add the stopContainerCmd
func init() {
	rootCmd.AddCommand(stopContainerCmd)
	stopContainerCmd.Flags().StringVarP(&containerID, "id", "i", "", "Container ID to use for stopping the container")
}

// Add the stopContainer function
func stopContainer(containerIDFlag string) {
	if containerIDFlag == "" {
		var err error
		containerID, err = helpers.GetContainerID()
		if err != nil {
			fmt.Println("Error reading container ID:", err)
			return
		}
	} else {
		containerID = containerIDFlag
	}

	apiUrl := fmt.Sprintf("https://localhost:8440/api/v1/containers/%s/stop", containerID)
	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

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

	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("Container stopped successfully.")
	} else {
		fmt.Println("Error stopping container:", resp.Status)
	}
}

//./borrower-cli stop-container --id 99d69fd86a9eda71ff46793af442031f303e246cdbbfcf226bbc3fab09a31768
//./borrower-cli stop-container
