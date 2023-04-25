package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
)

// Resource represents the resource structure to send to the Resource Manager
type Resource struct {
	ResourceType string `json:"type"`
	CPUCores     int    `json:"cpu_cores"`
	MemoryMB     int    `json:"memory_mb"`
	StorageGB    int    `json:"storage_gb"`
	LenderID     string `json:"lender_id,omitempty"`
}

func init() {
	lendCmd.Flags().String("resource-type", "", "Resource type (required)")
	lendCmd.Flags().Int("cpu-cores", 0, "CPU cores to request")
	lendCmd.Flags().Int("memory-mb", 0, "Memory in MB to request")
	lendCmd.Flags().Int("storage-gb", 0, "Storage in GB to request")
	rootCmd.AddCommand(lendCmd)
}

var lendCmd = &cobra.Command{
	Use:   "lend",
	Short: "Lend resources to the Resource Manager",
	Long:  `This command allows lenders to specify resources they want to lend, such as CPU, RAM, and storage.`,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType, _ := cmd.Flags().GetString("resource-type")
		cpuCores, _ := cmd.Flags().GetInt("cpu-cores")
		memoryMB, _ := cmd.Flags().GetInt("memory-mb")
		storageGB, _ := cmd.Flags().GetInt("storage-gb")

		// Create a resource instance
		resource := Resource{
			ResourceType: resourceType,
			CPUCores:     cpuCores,
			MemoryMB:     memoryMB,
			StorageGB:    storageGB,
			//LenderID: "",
		}

		fmt.Printf("Lending resources: Type: %s, CPU Cores: %d, Memory: %d MB, Storage: %d GB\n", resourceType, cpuCores, memoryMB, storageGB)

		// Send the resource data to the Resource Manager
		registerResource(resource)
	},
}

func registerResource(resource Resource) {
	apiGatewayURL := "http://localhost:8081/api/v1/resources"

	// Marshal the resource data into JSON
	jsonData, err := json.Marshal(resource)
	if err != nil {
		log.Fatalf("Error marshalling resource data: %v", err)
	}

	// Send an HTTP POST request to the Resource Manager with the JSON data
	resp, err := http.Post(apiGatewayURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error sending HTTP request to Resource Manager: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf("Resource registered successfully: %s\n", bodyString)
	} else {
		log.Fatalf("Error registering resource: status code %d", resp.StatusCode)
	}
}

//./lender-cli lend --resource-type vm --cpu-cores 2 --memory-mb 1024 --storage-gb 2
