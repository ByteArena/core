package utils

var (
	version string
)

func GetVersion() string {

	if version == "" {
		return "dev (unspecified)"
	}

	return version
}
