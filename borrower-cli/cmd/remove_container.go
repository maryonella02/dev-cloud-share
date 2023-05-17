package cmd

import (
	"borrower-cli/helpers"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
)

// Add a new removeContainerCmd
var removeContainerCmd = &cobra.Command{
	Use:   "remove-container",
	Short: "Remove a container",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Removing container...")
		removeContainer(containerID)
	},
}

// Modify the init function to add the removeContainerCmd
func init() {
	rootCmd.AddCommand(removeContainerCmd)
	removeContainerCmd.Flags().StringVarP(&containerID, "id", "i", "", "Container ID to use for removing the container")
}

// Add the removeContainerCmd function
func removeContainer(containerIDFlag string) {
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

	apiUrl := fmt.Sprintf("https://localhost:8440/api/v1/containers/%s/remove", containerID)
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
		fmt.Println("Container removed successfully.")
	} else {
		fmt.Println("Error removing container:", resp.Status)
	}
}

//./borrower-cli remove-container --id 99d69fd86a9eda71ff46793af442031f303e246cdbbfcf226bbc3fab09a31768
//./borrower-cli remove-container
