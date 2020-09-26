package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/mattermost/mattermost-server/v5/model"
)

func getClient() (*model.Client4, error) {
	socketPath := os.Getenv("MM_LOCALSOCKETPATH")
	if socketPath == "" {
		socketPath = model.LOCAL_MODE_SOCKET_PATH
	}

	client, connected := getUnixClient(socketPath)
	if connected {
		log.Printf("Connecting using local mode over %s", socketPath)
		return client, nil
	}

	if os.Getenv("MM_LOCALSOCKETPATH") != "" {
		log.Printf("No socket found at %s for local mode deployment. Attempting to authenticate with credentials.", socketPath)
	}

	siteURL := os.Getenv("MM_SERVICESETTINGS_SITEURL")
	adminToken := os.Getenv("MM_ADMIN_TOKEN")
	adminUsername := os.Getenv("MM_ADMIN_USERNAME")
	adminPassword := os.Getenv("MM_ADMIN_PASSWORD")

	if siteURL == "" {
		return nil, errors.New("MM_SERVICESETTINGS_SITEURL is not set")
	}

	client = model.NewAPIv4Client(siteURL)

	if adminToken != "" {
		log.Printf("Authenticating using token against %s.", siteURL)
		client.SetToken(adminToken)
		return client, nil
	}

	if adminUsername != "" && adminPassword != "" {
		client := model.NewAPIv4Client(siteURL)
		log.Printf("Authenticating as %s against %s.", adminUsername, siteURL)
		_, resp := client.Login(adminUsername, adminPassword)
		if resp.Error != nil {
			return nil, fmt.Errorf("failed to login as %s: %w", adminUsername, resp.Error)
		}
		return client, nil
	}

	return nil, errors.New("one of MM_ADMIN_TOKEN or MM_ADMIN_USERNAME/MM_ADMIN_PASSWORD must be defined")
}

func getUnixClient(socketPath string) (*model.Client4, bool) {
	_, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, false
	}

	return model.NewAPIv4SocketClient(socketPath), true
}
