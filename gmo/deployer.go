package gmo

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type deployment struct {
	containers  map[string]snakeContainer
	cli         *client.Client
	defaultPort int
}

type snakeContainer struct {
	envVariables  []string
	containerId   string
	imageName     string
	containerName string
	port          string
}

func createDeployment() *deployment {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	deployment := &deployment{
		cli:         cli,
		defaultPort: 8081,
		containers:  make(map[string]snakeContainer),
	}

	return deployment
}

func (d *deployment) run(ctx context.Context) bool {
	// Create the containers
	d.createContainers(ctx)

	// Run the containers
	d.startContainers(ctx)

	// Health check
	isSuccess := d.runHealthCheck()

	if !isSuccess {
		d.stopAndRemoveContainers(ctx)
	}

	return isSuccess
}

func (d *deployment) runHealthCheck() bool {
	time.Sleep(5 * time.Second)
	for name, con := range d.containers {
		url := fmt.Sprintf("http://localhost:%s", con.port)

		log.Println("RUNNING HEALTH CHECK", name, url)
		res, err := http.Get(url)
		if err != nil {
			log.Println(err)
			return false
		}

		if res.StatusCode != 200 {
			log.Println("Container did not respond with 200 during health check", name, con.imageName, con.port)
			return false
		}
		log.Println("SUCCESS")

		res.Body.Close()
	}

	return true
}

func (d *deployment) pullImages(ctx context.Context) {
	log.Println("PULLING IMAGES")
	for _, con := range d.containers {
		reader, err := d.cli.ImagePull(ctx, con.imageName, types.ImagePullOptions{})
		if err != nil {
			log.Fatal(err)
		}

		defer reader.Close()
		io.Copy(os.Stdout, reader)
	}
}

func (d *deployment) createContainers(ctx context.Context) {
	log.Println("CREATING CONTAINERS")
	for name, con := range d.containers {
		config := &container.Config{
			Image: con.imageName,
			ExposedPorts: nat.PortSet{
				nat.Port("8080/tcp"): struct{}{},
			},
			Env: con.envVariables,
		}

		hostConfig := &container.HostConfig{
			PortBindings: nat.PortMap{
				nat.Port("8080/tcp"): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: con.port,
					},
				},
			},
		}

		resp, err := d.cli.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
		if err != nil {
			log.Fatal(err)
		}

		con.containerId = resp.ID

		d.containers[name] = con
	}
}

func (d *deployment) startContainers(ctx context.Context) {
	log.Println("STARTING CONTAINERS")
	for _, con := range d.containers {
		if err := d.cli.ContainerStart(ctx, con.containerId, types.ContainerStartOptions{}); err != nil {
			log.Fatal(err)
		}
	}
}

func (d *deployment) stopAndRemoveContainers(ctx context.Context) {
	log.Println("STOPPING AND REMOVING CONTAINERS")
	defer d.cli.Close()
	for _, con := range d.containers {
		if err := d.cli.ContainerStop(ctx, con.containerId, nil); err != nil {
			log.Fatal(err)
		}

		if err := d.cli.ContainerRemove(ctx, con.containerId, types.ContainerRemoveOptions{}); err != nil {
			log.Fatal(err)
		}
	}
}

func (d *deployment) addContainer(name, imageName string, envVariables []string) {
	c := snakeContainer{
		imageName:     imageName,
		containerName: name,
		port:          fmt.Sprintf("%d", d.defaultPort),
		envVariables:  envVariables,
	}
	d.defaultPort += 1
	d.containers[name] = c
}