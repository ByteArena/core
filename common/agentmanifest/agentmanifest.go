package agentmanifest

import (
	"encoding/json"
	"io/ioutil"
	"path"

	arenaservertypes "github.com/bytearena/core/arenaserver/types"
	bettererrors "github.com/xtuc/better-errors"
)

type AgentManifest struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Author      string `json:"author"`
	License     string `json:"license"`
	Language    string `json:"language"`
	GameMode    string `json:"gameMode"`
	RepoURL     string `json:"RepoURL"`
	Description string `json:"description"`
	AvatarURL   string `json:"avatarURL"`
}

const (
	AGENT_MANIFEST_LABEL_KEY = "bytearena.manifest"
	AGENT_MANIFEST_FILENAME  = "ba.json"
)

func GetByAgentContainer(
	container *arenaservertypes.AgentContainer,
	orch arenaservertypes.ContainerOrchestrator,
) (*AgentManifest, error) {

	inspectResult, _, _ := orch.GetCli().ImageInspectWithRaw(
		orch.GetContext(),
		container.ImageName,
	)

	labels := inspectResult.Config.Labels

	agentManifest, err := ParseFromString(
		[]byte(labels[AGENT_MANIFEST_LABEL_KEY]),
	)

	if err != nil {
		return err
	}

	return agentManifest
}

func ParseFromString(content []byte) (*AgentManifest, error) {
	var manifest AgentManifest

	err := json.Unmarshal(content, &manifest)

	return &manifest, err
}

func ParseFromDir(dir string) (*AgentManifest, error) {
	fileLocation := path.Join(dir, AGENT_MANIFEST_FILENAME)

	content, err := ioutil.ReadFile(fileLocation)

	if err != nil {
		return nil, bettererrors.
			New("Parsing agent's manifest error").
			SetContext("filename", AGENT_MANIFEST_FILENAME).
			With(bettererrors.NewFromErr(err))
	}

	return ParseFromString(content)
}

func Validate(manifest *AgentManifest) error {

	if manifest == nil {
		return bettererrors.New("Invalid manifest file")
	}

	if manifest.Id == "" {
		return bettererrors.New("Missing id")
	}

	if manifest.Name == "" {
		return bettererrors.New("Missing name")
	}

	return nil
}

func (manifest AgentManifest) ToString() string {
	bytes, _ := json.Marshal(manifest)

	return string(bytes)
}
