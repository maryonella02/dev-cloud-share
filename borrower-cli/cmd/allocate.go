package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
)

type ResourceRequest struct {
	ResourceType string `json:"resource_type"`
	MinCPUCores  int    `json:"min_cpu_cores"`
	MinMemoryMB  int    `json:"min_memory_mb"`
	MinStorageGB int    `json:"min_storage_gb"`
}

func init() {
	allocateCmd.Flags().String("resource-type", "", "Resource type")
	allocateCmd.Flags().Int("min-cpu-cores", 0, "Minimum number of CPU cores")
	allocateCmd.Flags().Int("min-memory-mb", 0, "Minimum memory in MB")
	allocateCmd.Flags().Int("min-storage-gb", 0, "Minimum storage in GB")
	rootCmd.AddCommand(allocateCmd)
}

var allocateCmd = &cobra.Command{
	Use:   "allocate",
	Short: "Allocate a resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		resourceType, _ := cmd.Flags().GetString("resource-type")
		minCPUCores, _ := cmd.Flags().GetInt("min-cpu-cores")
		minMemoryMB, _ := cmd.Flags().GetInt("min-memory-mb")
		minStorageGB, _ := cmd.Flags().GetInt("min-storage-gb")

		var allocationInfo struct {
			BorrowerID      string          `json:"borrower_id"`
			ResourceRequest ResourceRequest `json:"resource_request"`
		}
		//allocationInfo.BorrowerID = "1"
		// TODO: add normal borrowerID identification and allocation
		allocationInfo.ResourceRequest = ResourceRequest{
			ResourceType: resourceType,
			MinCPUCores:  minCPUCores,
			MinMemoryMB:  minMemoryMB,
			MinStorageGB: minStorageGB,
		}

		// Replace this URL with the API Gateway URL for the allocations endpoint
		url := "http://localhost:8081/api/v1/allocations"

		client := &http.Client{}
		requestJSON, _ := json.Marshal(allocationInfo)
		req, err := http.NewRequest("POST", url, bytes.NewReader(requestJSON))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(string(body))
			fmt.Println(err)
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to allocate resource: %s", string(body))
		}

		fmt.Printf("Resource allocated successfully: %s\n", string(body))
		return nil
	},
}
