package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

func init() {
	rootCmd.AddCommand(requestCmd)

	requestCmd.Flags().String("resource-type", "", "Resource type (required)")
	requestCmd.Flags().Int("min-cpu-cores", 0, "Minimum number of CPU cores to request")
	requestCmd.Flags().Int("min-memory-mb", 0, "Minimum memory in MB to request")
	requestCmd.Flags().Int("min-storage-gb", 0, "Minimum storage in GB to request")

	// Mark the "type" flag as required
	requestCmd.MarkFlagRequired("resource-type")

	// Bind flags to Viper configuration
	viper.BindPFlag("resource-type", requestCmd.Flags().Lookup("resource-type"))
	viper.BindPFlag("min-cpu-cores", requestCmd.Flags().Lookup("min-cpu-cores"))
	viper.BindPFlag("min-memory-mb", requestCmd.Flags().Lookup("min-memory-mb"))
	viper.BindPFlag("min-storage-gb", requestCmd.Flags().Lookup("min-storage-gb"))
}

type ResourceRequest struct {
	ResourceType string `json:"resource_type"`
	MinCPUCores  int    `json:"min_cpu_cores"`
	MinMemoryMB  int    `json:"min_memory_mb"`
	MinStorageGB int    `json:"min_storage_gb"`
}

const APIAllocationsURL = "http://localhost:8081/api/v1/allocations"

var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "Request resources from the Resource Manager",
	Long:  `This command allows borrowers to specify resources they want to request, such as CPU, RAM, and storage.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		resourceType, _ := cmd.Flags().GetString("resource-type")
		minCPUCores, _ := cmd.Flags().GetInt("min-cpu-cores")
		minMemoryMB, _ := cmd.Flags().GetInt("min-memory-mb")
		minStorageGB, _ := cmd.Flags().GetInt("min-storage-gb")

		var allocationInfo struct {
			BorrowerID      string          `json:"borrower_id"`
			ResourceRequest ResourceRequest `json:"resource_request"`
		}

		allocationInfo.BorrowerID = "6446df1322b3d57d49cc2264"
		// TODO: add normal borrowerID identification and allocation
		allocationInfo.ResourceRequest = ResourceRequest{
			ResourceType: resourceType,
			MinCPUCores:  minCPUCores,
			MinMemoryMB:  minMemoryMB,
			MinStorageGB: minStorageGB,
		}

		fmt.Printf("Requesting resources: Type: %s, CPU Cores: %d, Memory: %d MB, Storage: %d GB\n", resourceType, minCPUCores, minMemoryMB, minStorageGB)

		client := &http.Client{}
		requestJSON, _ := json.Marshal(allocationInfo)
		req, err := http.NewRequest("POST", APIAllocationsURL, bytes.NewReader(requestJSON))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(string(body))
			fmt.Println(err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to allocate resource: %s for body %s", string(body), string(requestJSON))
		}

		fmt.Printf("Resource allocated successfully!")
		return nil
	},
}

// ./borrower-cli request --resource-type example_type --min-cpu-cores 4 --min-memory-mb 8192 --min-storage-gb 100
