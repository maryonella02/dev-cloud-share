package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/term"
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
	Short: "Login as a borrower",
	Run: func(cmd *cobra.Command, args []string) {
		username := getUsername()
		password := getPassword()

		user := User{
			Username: username,
			Password: password,
		}

		fmt.Printf("Logging in as borrower: Username: %s\n", username)

		Token = loginUser(user)
		err := os.WriteFile("token.txt", []byte(Token), 0644)
		if err != nil {
			fmt.Println("Error saving token:", err)
		} else {
			fmt.Println("Token saved to token.txt")
		}
		fmt.Printf("Login successful. Token: %s\n", Token)
	},
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register as a borrower",
	Run: func(cmd *cobra.Command, args []string) {
		username := getUsername()
		password := getPassword()
		role := getRole()

		user := User{
			Username: username,
			Password: password,
			Role:     role,
		}

		fmt.Printf("Registering as borrower: Username: %s", username)

		Token = registerUser(user)
		fmt.Printf("Registration successful. Token: %s\n", Token)
	},
}

type Response struct {
	Token string `json:"token"`
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
		var response Response

		err = json.NewDecoder(resp.Body).Decode(&response)
		return response.Token
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
		var response Response

		err = json.NewDecoder(resp.Body).Decode(&response)
		return response.Token
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
	return "borrower"
}
