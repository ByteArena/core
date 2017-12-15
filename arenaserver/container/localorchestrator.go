package container

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	uuid "github.com/satori/go.uuid"

	"github.com/bytearena/core/common/types"
	"github.com/bytearena/core/common/utils"
)

type LocalContainerOrchestrator struct {
	ctx          context.Context
	cli          *client.Client
	registryAuth string
	host         string
	containers   []*types.AgentContainer
	events       chan interface{}
}

const (
	LOG_ENTRY_BUFFER = 100
)

func (orch *LocalContainerOrchestrator) startContainerLocalOrch(ctner *types.AgentContainer, addTearDownCall func(types.TearDownCallback)) error {

	err := orch.cli.ContainerStart(
		orch.ctx,
		ctner.Containerid,
		dockertypes.ContainerStartOptions{},
	)

	if err != nil {
		return err
	}

	err = orch.localLogsToStdOut(ctner)

	if err != nil {
		return errors.New("Failed to follow docker container logs for " + ctner.Containerid)
	}

	containerInfo, err := orch.cli.ContainerInspect(
		orch.ctx,
		ctner.Containerid,
	)
	if err != nil {
		return errors.New("Could not inspect container " + ctner.Containerid)
	}

	ctner.SetIPAddress(containerInfo.NetworkSettings.IPAddress)

	return nil
}

func MakeLocalContainerOrchestrator(host string) types.ContainerOrchestrator {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	utils.Check(err, "Failed to initialize docker client environment")

	registryAuth := ""

	return &LocalContainerOrchestrator{
		ctx:          ctx,
		cli:          cli,
		host:         host,
		registryAuth: registryAuth,
		events:       make(chan interface{}, LOG_ENTRY_BUFFER),
	}
}

func (orch *LocalContainerOrchestrator) GetHost() (string, error) {
	if orch.host == "" {
		res, err := orch.cli.NetworkInspect(orch.ctx, "bridge", dockertypes.NetworkInspectOptions{})
		if err != nil {
			return "", err
		}

		return res.IPAM.Config[0].Gateway, nil
	}

	return orch.host, nil
}

func (orch *LocalContainerOrchestrator) localLogsToStdOut(container *types.AgentContainer) error {

	go func(orch *LocalContainerOrchestrator, container *types.AgentContainer) {

		reader, err := orch.cli.ContainerLogs(orch.ctx, container.Containerid, dockertypes.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
			Details:    false,
			Timestamps: false,
		})

		utils.Check(err, "Could not read container logs for "+container.AgentId.String()+"; container="+container.Containerid)

		defer reader.Close()

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			buf := scanner.Bytes()

			/*
				This is to remove Docker log header.
				First 8 bytes are part of the header.
			*/
			if len(buf) > 8 {
				buf = buf[8:]
			}

			orch.events <- EventAgentLog{
				Value:     string(buf),
				AgentName: container.ImageName,
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

	}(orch, container)

	return nil
}

func (orch *LocalContainerOrchestrator) StartAgentContainer(ctner *types.AgentContainer, addTearDownCall func(types.TearDownCallback)) error {
	orch.events <- EventDebug{"Spawning agent " + ctner.ImageName}

	return orch.startContainerLocalOrch(ctner, addTearDownCall)
}

func (orch *LocalContainerOrchestrator) CreateAgentContainer(agentid uuid.UUID, host string, port int, dockerimage string) (*types.AgentContainer, error) {
	return CommonCreateAgentContainer(orch, agentid, host, port, dockerimage)
}

func (orch *LocalContainerOrchestrator) TearDown(container *types.AgentContainer) {
	timeout := 5 * time.Second

	err := orch.cli.ContainerStop(
		orch.ctx,
		container.Containerid,
		&timeout,
	)

	if err != nil {
		orch.events <- EventDebug{"Killing container " + container.Containerid}
		orch.cli.ContainerKill(orch.ctx, container.Containerid, "KILL")
	}
}

func (orch *LocalContainerOrchestrator) RemoveAgentContainer(ctner *types.AgentContainer) error {

	// We don't want to remove images in local mode
	return nil
}

func (orch *LocalContainerOrchestrator) Wait(ctner *types.AgentContainer) (<-chan container.ContainerWaitOKBody, <-chan error) {
	waitChan, errorChan := orch.cli.ContainerWait(
		orch.ctx,
		ctner.Containerid,
		container.WaitConditionRemoved,
	)

	return waitChan, errorChan
}

func (orch *LocalContainerOrchestrator) SetAgentLogger(container *types.AgentContainer) error {
	// TODO(sven): implement log to stdout here
	return nil
}

func (orch *LocalContainerOrchestrator) TearDownAll() error {
	for _, container := range orch.containers {
		orch.TearDown(container)
	}

	return nil
}

func (orch *LocalContainerOrchestrator) GetCli() *client.Client {
	return orch.cli
}

func (orch *LocalContainerOrchestrator) GetContext() context.Context {
	return orch.ctx
}

func (orch *LocalContainerOrchestrator) GetRegistryAuth() string {
	return orch.registryAuth
}

func (orch *LocalContainerOrchestrator) AddContainer(ctner *types.AgentContainer) {
	orch.containers = append(orch.containers, ctner)
}

func (orch *LocalContainerOrchestrator) RemoveContainer(ctner *types.AgentContainer) {
	containers := make([]*types.AgentContainer, 0)

	for _, c := range orch.containers {
		if c.AgentId != ctner.AgentId {
			containers = append(containers, c)
		}
	}

	orch.containers = containers
}

func (orch *LocalContainerOrchestrator) Events() chan interface{} {
	return orch.events
}
