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
	Email    string `json:"email"`
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
		email := getEmail()
		password := getPassword()

		user := User{
			Email:    email,
			Password: password,
		}

		fmt.Printf("Logging in as borrower: Email: %s\n", email)

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
		email := getEmail()
		password := getPassword()
		role := getRole()

		user := User{
			Email:    email,
			Password: password,
			Role:     role,
		}

		fmt.Printf("Registering as borrower: Email: %s", email)

		Token = registerUser(user)
		if Token != "" {
			fmt.Printf("Registration successful. Token: %s\n", Token)

		} else {
			fmt.Printf("Registration failed.")
		}
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
func getEmail() string {
	fmt.Print("Enter email: ")
	var email string
	_, err := fmt.Scanln(&email)
	if err != nil {
		fmt.Println("Error reading email:", err)
		os.Exit(1)
	}
	return strings.TrimSpace(email)
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
