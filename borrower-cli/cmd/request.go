package cmd

import (
	"borrower-cli/helpers"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

func init() {
	rootCmd.AddCommand(requestCmd)

	requestCmd.Flags().Int("max-cpu-cores", 0, "Maximum number of CPU cores to request")
	requestCmd.Flags().Int("max-memory-mb", 0, "Maximum memory in MB to request")

	// Bind flags to Viper configuration
	viper.BindPFlag("max-cpu-cores", requestCmd.Flags().Lookup("max-cpu-cores"))
	viper.BindPFlag("max-memory-mb", requestCmd.Flags().Lookup("max-memory-mb"))
}

type Resource struct {
	CPUCores int `json:"cpu_cores"`
	MemoryMB int `json:"memory_mb"`
}

const APIAllocationsURL = "https://localhost:8440/api/v1/allocations"

var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "Request resources from the Resource Manager",
	Long:  `This command allows borrowers to specify resources they want to request, such as CPU, RAM, and storage.`,
	Run: func(cmd *cobra.Command, args []string) {
		minCPUCores, _ := cmd.Flags().GetInt("max-cpu-cores")
		minMemoryMB, _ := cmd.Flags().GetInt("max-memory-mb")

		var allocationInfo struct {
			BorrowerID      string   `json:"borrower_id"`
			ResourceRequest Resource `json:"resource"`
		}

		allocationInfo.BorrowerID = "6446df1322b3d57d49cc2264"
		// TODO: add normal borrowerID identification and allocation
		allocationInfo.ResourceRequest = Resource{
			CPUCores: minCPUCores,
			MemoryMB: minMemoryMB,
		}

		fmt.Printf("Requesting resource with: CPU Cores: %d, Memory: %d MB", minCPUCores, minMemoryMB)

		// Create a custom HTTP client with insecure TLS configuration
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		requestJSON, _ := json.Marshal(allocationInfo)
		req, err := http.NewRequest("POST", APIAllocationsURL, bytes.NewReader(requestJSON))
		if err != nil {
			fmt.Printf("\nError requesting resource: %v\n", err)
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
			fmt.Printf("\nError requesting resource: %v\n", err)
		}

		_, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("\nError requesting resource: %v\n", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("\nError requesting resource: %v status \n", resp.StatusCode)
		} else {
			fmt.Println("\nResource allocated successfully!")
		}
	},
}

// ./borrower-cli request  --max-cpu-cores 4 --max-memory-mb 8192
