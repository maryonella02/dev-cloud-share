package cmd

import (
	"borrower-cli/helpers"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
)

var containerID string
var startContainerCmd = &cobra.Command{
	Use:   "start-container",
	Short: "Start a container",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting container...")
		startContainer(containerID)
	},
}

func init() {
	rootCmd.AddCommand(startContainerCmd)
	startContainerCmd.Flags().StringVarP(&containerID, "id", "i", "", "Container ID to use for starting the container")
}

func startContainer(containerIDFlag string) {
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

	apiUrl := fmt.Sprintf("https://localhost:8440/api/v1/containers/%s/start", containerID)
	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("Container started successfully.")
	} else {
		fmt.Println("Error starting container:", resp.Status)
	}
}

//./borrower-cli start-container --id 99d69fd86a9eda71ff46793af442031f303e246cdbbfcf226bbc3fab09a31768
//./borrower-cli start-container

// TODO: add container name flag to control container with name, not ID
