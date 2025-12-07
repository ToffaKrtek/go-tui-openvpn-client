package types

import (
	"strings"
)

type Item struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

func (i Item) Title() string       { return i.Name }
func (i Item) Description() string { return i.Text }
func (i Item) FilterValue() string { return i.Name }

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
