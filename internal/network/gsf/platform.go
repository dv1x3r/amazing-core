package gsf

import "strings"

type Platform uint8

const (
	PlatformUnknown Platform = iota
	PlatformWindows
	PlatformOSX
)

func (p Platform) String() string {
	switch p {
	case PlatformWindows:
		return "windows"
	case PlatformOSX:
		return "osx"
	default:
		return "unknown"
	}
}

func ParsePlatformFromMachineOS(value string) Platform {
	s := strings.ToLower(value)
	switch {
	case strings.HasPrefix(s, "windows"):
		return PlatformWindows
	case strings.HasPrefix(s, "mac os x"):
		return PlatformOSX
	default:
		return PlatformUnknown
	}
}
