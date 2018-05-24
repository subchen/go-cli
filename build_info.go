package cli

import (
	"regexp"
)

// BuildInfo stores app build info
type BuildInfo struct {
	Timestamp   string
	GitBranch   string
	GitCommit   string
	GitRevCount string
}

// ParseBuildInfo parse a buildinfo string info struct
func ParseBuildInfo(info string) *BuildInfo {
	return &BuildInfo{
		Timestamp:   _readValue(info, "time"),
		GitBranch:   _readValue(info, "branch"),
		GitCommit:   _readValue(info, "commit"),
		GitRevCount: _readValue(info, "patches"),
	}
}

func _readValue(input, name string) string {
	re := regexp.MustCompile(`(^|\s)` + name + `:("[^"]*"|'[^']*'|[[:graph:]]+)($|\s)`)
	matched := re.FindAllStringSubmatch(input, 1)
	if matched != nil {
		value := matched[0][2]
		if value[0] == '"' || value[0] == '\'' {
			value = value[1 : len(value)-1]
		}
		return value
	}
	return ""
}
