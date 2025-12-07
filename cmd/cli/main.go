package main

import (
	"flag"
	"fmt"

	"github.com/ToffaKrtek/go-tui-openvpn-client/internal/cmd"
)

func main() {
	activateName := flag.String("a", "", "activate config")
	disconnectName := flag.String("d", "", "disconnect config")
	listConfigs := flag.Bool("l", false, "list configs")
	listSessions := flag.Bool("s", false, "list sessions")
	importPath := flag.String("path", "", "path to import")
	importName := flag.String("name", "", "name to import")
	deleteConfig := flag.String("delete", "", "delete config")
	flag.Parse()

	if len(*deleteConfig) > 0 {
		cmd.DeleteConfig(*deleteConfig)
	}

	if len(*importPath) > 0 && len(*importName) > 0 {
		cmd.ImportConfig(*importName, *importPath)
	}
	if len(*activateName) > 0 {
		cmd.ActiveConfig(*activateName)
	}
	if len(*disconnectName) > 0 {
		cmd.DisconnectSession(*disconnectName)
	}
	if *listConfigs {
		items, err := cmd.GetConfigs()
		if err == nil {
			fmt.Println(items)
		}
	}
	if *listSessions {
		items, err := cmd.GetSession()
		if err == nil {
			fmt.Println(items)
		}
	}
}
