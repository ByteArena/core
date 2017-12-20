package types

import (
	"encoding/json"
	"io/ioutil"
	"path"

	bettererrors "github.com/xtuc/better-errors"
)

type AgentManifest struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Author      string `json:"author"`
	License     string `json:"license"`
	Language    string `json:"language"`
	GameMode    string `json:"gamemode"`
	RepoURL     string `json:"repourl"`
	Description string `json:"description"`
	AvatarURL   string `json:"avatarurl"`
}

const (
	AGENT_MANIFEST_LABEL_KEY = "bytearena.manifest"
	AGENT_MANIFEST_FILENAME  = "ba.json"
)

func GetAgentManifestByDockerImageName(
	dockerImageName string,
	orch ContainerOrchestrator,
) (AgentManifest, error) {

	inspectResult, _, inspectResulterr := orch.GetCli().ImageInspectWithRaw(
		orch.GetContext(),
		dockerImageName,
	)

	if inspectResulterr != nil {
		return AgentManifest{}, bettererrors.
			New("Could not find Docker image").
			SetContext("image", dockerImageName)
	}

	labels := inspectResult.Config.Labels

	manifestString := labels[AGENT_MANIFEST_LABEL_KEY]

	if manifestString == "" {
		return AgentManifest{}, bettererrors.
			New("Manifest not found, are you sure it is an agent?").
			SetContext("image", dockerImageName)
	}

	agentManifest, err := ParseAgentManifestFromString(
		[]byte(manifestString),
	)

	if err != nil {
		return agentManifest, err
	}

	return agentManifest, nil
}

func ParseAgentManifestFromString(content []byte) (AgentManifest, error) {
	var manifest AgentManifest

	err := json.Unmarshal(content, &manifest)

	return manifest, err
}

func ParseAgentManifestFromDir(dir string) (AgentManifest, error) {
	fileLocation := path.Join(dir, AGENT_MANIFEST_FILENAME)

	content, err := ioutil.ReadFile(fileLocation)

	if err != nil {
		return AgentManifest{}, bettererrors.
			New("Parsing agent's manifest error").
			SetContext("filename", AGENT_MANIFEST_FILENAME).
			With(bettererrors.NewFromErr(err))
	}

	return ParseAgentManifestFromString(content)
}

func ValidateAgentManifest(manifest AgentManifest) error {

	if manifest.Id == "" {
		return bettererrors.New("Missing id")
	}

	if manifest.Name == "" {
		return bettererrors.New("Missing name")
	}

	return nil
}

func (manifest AgentManifest) String() string {
	bytes, _ := json.Marshal(manifest)

	return string(bytes)
}
