package helpers

import (
	"os"
)

func GetContainerID() (string, error) {
	containerID, err := os.ReadFile("container_id.txt")
	if err != nil {
		return "", err
	}
	return string(containerID), nil
}
