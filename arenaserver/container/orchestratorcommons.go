package container

import (
	"errors"
	"os"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/pkg/term"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/jsonmessage"
	uuid "github.com/satori/go.uuid"
	bettererrors "github.com/xtuc/better-errors"

	arenaservertypes "github.com/bytearena/core/arenaserver/types"
	"github.com/bytearena/core/common/agentmanifest"
	"github.com/bytearena/core/common/utils"
)

func normalizeDockerRef(dockerimage string) (string, error) {

	p, _ := reference.Parse(dockerimage)
	named, ok := p.(reference.Named)
	if !ok {
		return "", errors.New("Invalid docker image name")
	}

	parsedRefWithTag := reference.TagNameOnly(named)
	return parsedRefWithTag.String(), nil
}

func CommonCreateAgentContainer(orch arenaservertypes.ContainerOrchestrator, agentid uuid.UUID, host string, port int, dockerimage string) (*arenaservertypes.AgentContainer, error) {
	containerUnixUser := utils.GetenvOrDefault("CONTAINER_UNIX_USER", "root")

	normalizedDockerimage, err := normalizeDockerRef(dockerimage)

	if err != nil {
		return nil, bettererrors.NewFromErr(err)
	}

	localimages, _ := orch.GetCli().ImageList(orch.GetContext(), types.ImageListOptions{})
	foundlocal := false
	for _, localimage := range localimages {
		for _, alias := range localimage.RepoTags {
			if normalizedAlias, err := normalizeDockerRef(alias); err == nil {
				if normalizedAlias == normalizedDockerimage {
					foundlocal = true
					break
				}
			}
		}

		if foundlocal {
			break
		}
	}

	if !foundlocal {
		reader, err := orch.GetCli().ImagePull(
			orch.GetContext(),
			dockerimage,
			types.ImagePullOptions{
				RegistryAuth: orch.GetRegistryAuth(),
			},
		)

		if err != nil {
			return nil, bettererrors.
				New("Failed to pull from registry").
				With(bettererrors.NewFromErr(err)).
				SetContext("image", dockerimage)
		}

		fd, isTerminal := term.GetFdInfo(os.Stdout)

		if err := jsonmessage.DisplayJSONMessagesStream(reader, os.Stdout, fd, isTerminal, nil); err != nil {
			return nil, err
		}

		reader.Close()
	}

	// Remove this
	inspectResult, _, _ := orch.GetCli().ImageInspectWithRaw(
		orch.GetContext(),
		normalizedDockerimage,
	)

	labels := inspectResult.Config.Labels

	agentManifest, _ := agentmanifest.ParseFromString(
		[]byte(labels[agentmanifest.AGENT_MANIFEST_LABEL_KEY]),
	)

	spew.Dump(agentManifest)

	// Remove this

	containerconfig := container.Config{
		Image: normalizedDockerimage,
		User:  containerUnixUser,
		Env: []string{
			"PORT=" + strconv.Itoa(port),
			"HOST=" + host,
			"AGENTID=" + agentid.String(),
		},
		AttachStdout: false,
		AttachStderr: false,
	}

	hostconfig := container.HostConfig{
		CapDrop:        []string{"ALL"},
		Privileged:     false,
		AutoRemove:     true,
		ReadonlyRootfs: true,
		NetworkMode:    "bridge",
		// Resources: container.Resources{
		// 	Memory: 1024 * 1024 * 32, // 32M
		// 	//CPUQuota: 5 * (1000),       // 5% en cent-millièmes
		// 	//CPUShares: 1,
		// 	CPUPercent: 5,
		// },
	}

	resp, err := orch.GetCli().ContainerCreate(
		orch.GetContext(), // go context
		&containerconfig,  // container config
		&hostconfig,       // host config
		nil,               // network config
		"agent-"+agentid.String(), // container name
	)
	if err != nil {
		return nil, bettererrors.New("Failed to create docker container for agent " + agentid.String() + "; " + err.Error())
	}

	agentcontainer := arenaservertypes.NewAgentContainer(agentid, resp.ID, normalizedDockerimage)
	orch.AddContainer(agentcontainer)

	return agentcontainer, nil
}
