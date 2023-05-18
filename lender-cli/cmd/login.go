package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(registerCmd)
}

var Token string
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login as a lender",
	Run: func(cmd *cobra.Command, args []string) {
		username := getUsername()
		password := getPassword()

		user := User{
			Username: username,
			Password: password,
		}

		fmt.Printf("Logging in as lender: Username: %s\n", username)

		Token = loginUser(user)
		fmt.Printf("Login successful. Token: %s\n", Token)
	},
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register as a lender",
	Run: func(cmd *cobra.Command, args []string) {
		username := getUsername()
		password := getPassword()
		role := getRole()

		user := User{
			Username: username,
			Password: password,
			Role:     role,
		}

		fmt.Printf("Registering as lender: Username: %s", username)

		Token = registerUser(user)
		fmt.Printf("Registration successful. Token: %s\n", Token)
	},
}

func loginUser(user User) string {
	apiGatewayURL := "https://localhost:8440/api/v1/login"

	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("Error marshalling user data: %v", err)
	}

	resp, err := http.Post(apiGatewayURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error sending HTTP request to API Gateway: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}
		token := string(bodyBytes)
		return token
	} else {
		log.Fatalf("Error logging in: status code %d", resp.StatusCode)
	}

	return ""
}

func registerUser(user User) string {
	apiGatewayURL := "https://localhost:8440/api/v1/register"

	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("Error marshalling user data: %v", err)
	}

	resp, err := http.Post(apiGatewayURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error sending HTTP request to API Gateway: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}
		token := string(bodyBytes)
		return token
	} else {
		log.Fatalf("Error registering user: status code %d", resp.StatusCode)
	}

	return ""
}

// Helper functions to get input from the user
func getUsername() string {
	fmt.Print("Enter username: ")
	var username string
	_, err := fmt.Scanln(&username)
	if err != nil {
		fmt.Println("Error reading username:", err)
		os.Exit(1)
	}
	return strings.TrimSpace(username)
}

func getPassword() string {
	fmt.Print("Enter password: ")
	var password string
	// Use terminal input without echoing the characters
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Print a new line after reading the password
	if err != nil {
		fmt.Println("Error reading password:", err)
		os.Exit(1)
	}
	password = string(bytePassword)
	return strings.TrimSpace(password)
}

func getRole() string {
	return "lender"
}
