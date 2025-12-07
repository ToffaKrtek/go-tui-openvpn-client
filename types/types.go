package types

import (
	"strings"
)

type Item struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type Config struct {
	Name     string `json:"name"`
	LastUsed string `json:"lastused"`
}

func FindSessionName(data string) string {
	const prefix = "Config name:"
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if strings.Contains(line, prefix) {
			after := strings.TrimSpace(strings.SplitN(line, prefix, 2)[1])
			if idx := strings.IndexByte(after, ' '); idx != -1 {
				return after[:idx]
			}
			return after
		}
	}
	return ""
}
