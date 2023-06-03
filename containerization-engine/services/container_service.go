package services

import (
	"containerization-engine/models"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type ContainerService struct {
	cli *client.Client
}

func NewContainerService() (*ContainerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &ContainerService{cli: cli}, nil
}

func (cs *ContainerService) CreateContainer(config models.ContainerConfig) (string, error) {
	ctx := context.Background()

	// TODO: allow pulling private images. Use https://goharbor.io/

	_, err := cs.cli.ImagePull(ctx, config.Image, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}

	resp, err := cs.cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: config.Image,
			Cmd:   config.Command,
			Env:   convertEnvMapToSlice(config.Environment),
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory:   config.Memory,
				NanoCPUs: config.NanoCPUs,
			},
		},
		&network.NetworkingConfig{},
		nil, "",
	)

	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (cs *ContainerService) StartContainer(containerID string) error {
	ctx := context.Background()

	return cs.cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}

func (cs *ContainerService) StopContainer(containerID string) error {
	ctx := context.Background()

	return cs.cli.ContainerStop(ctx, containerID, container.StopOptions{})
}

func (cs *ContainerService) RemoveContainer(containerID string) error {
	ctx := context.Background()

	return cs.cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{})
}

func (cs *ContainerService) GetContainerStatus(containerID string) (string, error) {
	ctx := context.Background()

	containerInfo, err := cs.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", err
	}

	return containerInfo.State.Status, nil
}

func convertEnvMapToSlice(envMap map[string]string) []string {
	env := make([]string, 0, len(envMap))

	for k, v := range envMap {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	return env
}
