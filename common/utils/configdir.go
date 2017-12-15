package utils

import (
	"os/user"
	"path"
)

func getBAConfigDir() (string, error) {
	user, err := user.Current()

	if err != nil {
		return "", err
	}

	return path.Join(user.HomeDir, ".bytearena"), nil
}

func GetTrainerMapsDir() (string, error) {
	baConfigDir, err := getBAConfigDir()
	if err != nil {
		return "", err
	}

	return path.Join(baConfigDir, "maps"), nil
}

func GetTrainerMapsManifestPath() (string, error) {
	mapsDir, err := GetTrainerMapsDir()
	if err != nil {
		return "", err
	}

	return path.Join(mapsDir, "manifest.json"), nil
}
