package types

import "strings"

type FileLocationType int

const (
	FileLocationTypeLocal FileLocationType = iota
	FileLocationTypeRemoteHTTP
)

type ConfigFile string

func (config ConfigFile) String() string {
	return string(config)
}

func (config ConfigFile) LocationType() FileLocationType {
	if strings.HasPrefix(string(config), "http://") || strings.HasPrefix(string(config), "https://") {
		return FileLocationTypeRemoteHTTP
	}
	return FileLocationTypeLocal
}
