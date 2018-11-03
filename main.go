package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
	"os"
	"strings"
	"time"
)

// HostEntry ...
type HostEntry struct {
	ip   string
	host string
}

func (h *HostEntry) prepare() []byte {
	return []byte(fmt.Sprintf("%s %s\n", h.ip, h.host))
}

var (
	//ErrNoLabel blah blah
	ErrNoLabel = errors.New("error no label set")
)

func writeFile(hosts []HostEntry) {
	f, err := os.Create("/in/hosts")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, host := range hosts {
		f.Write(host.prepare())
	}
	f.Sync()
}

func getContainerIP(cli *client.Client, containerID string) ([]HostEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return []HostEntry{}, err
	}

	hosts := make([]HostEntry, 0)
	for _, container := range containers {
		if container.State != "running" {
			log.Printf("Skipping container %+v", container.Names)
			continue
		}

		if host, ok := container.Labels["dns.host"]; ok {
			var address string
			networks := container.NetworkSettings.Networks
			if value, ok := networks["bridge"]; ok {
				address = value.IPAddress
			} else {
			}
			for name, netConf := range networks {
				if name == "bridge" {
					address = netConf.IPAddress
					break
				}

				if strings.HasSuffix(name, "default") {
					address = netConf.IPAddress
					break
				}
			}

			if address == "" {
				log.Println("network has no default network, defaulting to 127.0.0.1")
				address = "127.0.0.1"
			}

			host := HostEntry{
				host: host,
				ip:   address,
			}
			hosts = append(hosts, host)
		}
	}

	return hosts, nil
}

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	events, chErr := cli.Events(context.Background(), types.EventsOptions{})

	log.Printf("listening to events")
	for {
		select {
		case event := <-events:
			if event.Type == "container" {
				attrs := event.Actor.Attributes
				domains, err := getContainerIP(cli, attrs["ID"])
				if err != nil {
					if err == ErrNoLabel {
						log.Printf("Label not set for container")
					} else {
						log.Fatal(err)
					}
				}
				go writeFile(domains)
			}

		case err := <-chErr:
			log.Fatal(err)
		}

	}
}
