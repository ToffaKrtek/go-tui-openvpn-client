package cmd

import (
	"encoding/json"
	"os/exec"

	"github.com/ToffaKrtek/go-tui-openvpn-client/types"
)

func GetConfigs() ([]types.Item, error) {
	var err error
	configs := []types.Item{}
	command := exec.Command("openvpn3", "configs-list", "--json")
	data, err := command.Output()
	if err == nil {
		var rawMap map[string]types.Config
		if err = json.Unmarshal(data, &rawMap); err == nil {
			for _, cfg := range rawMap {
				configs = append(configs, types.Item{
					Name: cfg.Name,
					Text: cfg.LastUsed,
				})
			}
		}
	}
	return configs, err
}

func GetSession() ([]types.Item, error) {
	var err error
	sessions := []types.Item{}
	command := exec.Command("openvpn3", "sessions-list")
	data, err := command.Output()
	if err == nil {
		activeSession := types.FindSessionName(string(data))
		if len(activeSession) > 0 {
			sessions = append(sessions, types.Item{Name: activeSession, Text: "Active"})
		}
	}
	return sessions, err
}

func ActiveConfig(name string) error {
	command := exec.Command("openvpn3", "session-start", "-c", name, "--background")
	_, err := command.Output()
	return err
}

func DisconnectSession(name string) error {
	command := exec.Command("openvpn3", "session-manage", "-c", name, "-D")
	_, err := command.Output()
	return err
}

func ImportConfig(name string, path string) error {
	command := exec.Command("openvpn3", "config-import", "-c", path, "--name", name, "--persistent")
	_, err := command.Output()
	return err
}

func DeleteConfig(name string) error {
	command := exec.Command("openvpn3", "config-remove", "-c", name, "--force")
	_, err := command.Output()
	return err
}
