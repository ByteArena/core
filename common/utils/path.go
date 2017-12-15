package utils

import (
	"path"

	"github.com/kardianos/osext"
)

func GetAbsoluteDir(relative string) string {

	exfolder, err := osext.ExecutableFolder()
	Check(err, "Cannot get absolute dir for "+relative)

	return path.Join(exfolder, relative)
}

func GetExecutablePath() string {
	binpath, _ := osext.Executable()
	return binpath
}

func GetExecutableDir() string {
	return path.Dir(GetExecutablePath())
}
