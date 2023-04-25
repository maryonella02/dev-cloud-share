package controllers

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"
	dockerContainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ContainerController struct {
	cli *client.Client
}

func NewContainerController() (*ContainerController, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &ContainerController{cli: cli}, nil
}

func (c *ContainerController) CreateContainer(image string, config *dockerContainer.Config, hostConfig *dockerContainer.HostConfig) (string, error) {
	ctx := context.Background()

	// Pull the image from Docker Hub if not available locally
	_, err := c.cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}

	// Create the dockerContainer
	resp, err := c.cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (c *ContainerController) StartContainer(containerID string) error {
	ctx := context.Background()
	return c.cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}

func (c *ContainerController) StopContainer(containerID string) error {
	ctx := context.Background()
	return c.cli.ContainerStop(ctx, containerID, dockerContainer.StopOptions{})
}

func (c *ContainerController) RemoveContainer(containerID string) error {
	ctx := context.Background()
	return c.cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{})
}

func (c *ContainerController) GetContainerStatus(containerID string) (string, error) {
	ctx := context.Background()
	container, err := c.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", err
	}

	if !container.State.Running && container.State.Error != "" {
		return container.State.Error, errors.New("dockerContainer error")
	}

	return container.State.Status, nil
}
