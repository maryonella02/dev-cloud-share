package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
)

func init() {
	releaseCmd.Flags().String("resource-id", "", "Resource ID to release (required)")
	rootCmd.AddCommand(releaseCmd)
}

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Release a requested resource",
	Long:  `This command allows borrowers to release a requested resource.`,
	Run: func(cmd *cobra.Command, args []string) {
		resourceID, _ := cmd.Flags().GetString("resource-id")

		err := releaseResource(resourceID)
		if err != nil {
			fmt.Printf("Error releasing resource: %v\n", err)
			return
		}
		fmt.Printf("Resource %s released successfully\n", resourceID)
	},
}

func releaseResource(resourceID string) error {
	apiURL := APIAllocationsURL + "/" + resourceID

	req, err := http.NewRequest("DELETE", apiURL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to release resource: %s", resp.Status)
	}

	return nil
}

// ./borrower-cli release --resource-id 6447f49c49c4ce579009297d

// TODO Add another command to get all resources in usage for borrower to get identifier from CLI
