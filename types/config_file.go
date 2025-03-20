package types

import "strings"

type FileLocationType int

const (
	FileLocationTypeLocal FileLocationType = iota
	FileLocationTypeRemoteHTTP
)

type ConfigFile string

func (configFile ConfigFile) String() string {
	return string(configFile)
}

func (configFile ConfigFile) LocationType() FileLocationType {
	if strings.HasPrefix(string(configFile), "http://") || strings.HasPrefix(string(configFile), "https://") {
		return FileLocationTypeRemoteHTTP
	}
	return FileLocationTypeLocal
}

func (configFile ConfigFile) Extension() string {
	i := strings.LastIndex(configFile.String(), ".")
	if i == -1 {
		return ""
	}

	return configFile.String()[i+1:]
}
