package dckr

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

// Config configuration for building docker clients
type Config struct {
	UseBoot2Docker bool
	Endpoint       string
	Description    string
}

func (c Config) String() string {
	if c.Endpoint != "" {
		return fmt.Sprintf("Endpoint: %s -- %s", c.Endpoint, c.Description)
	} else if c.UseBoot2Docker {
		return fmt.Sprintf("Boot2Docker[env] -- %s", c.Description)
	}
	return fmt.Sprintf("unknown -- %s", c.Description)
}

// Client a Docker client with easy to use API
type Client struct {
	cfg Config
	c   *docker.Client
}

// ImageID is a type for docker image ids
type ImageID string

// ContainerID is a type for docker image ids
type ContainerID string

// Image is a struct for Docker images
type Image struct {
	ID     ImageID
	Labels map[string]string
}

// Container is a struct for Docker containers
type Container struct {
	ID ContainerID
}

// NewClientFirstOf returns a new docker client or an error
func NewClientFirstOf(cfg []Config) (Client, error) {
	var buf bytes.Buffer
	for _, dcfg := range cfg {
		client, err := NewClient(dcfg)
		if err != nil || client.c == nil {
			buf.WriteString(fmt.Sprintf(
				"%s:\n\t%s\nReason:%s\n",
				"Failed to make new docker client using configuration",
				dcfg, err))
		} else if client.c != nil {
			log.Println("Successfully created Docker client using config: ", dcfg)
			return client, nil
		}
	}
	return Client{}, errors.New(buf.String())
}

// NewClient returns a new docker client or an error
func NewClient(dcfg Config) (Client, error) {
	var client *docker.Client
	var err error
	if dcfg.Endpoint != "" {
		client, err = docker.NewClient(dcfg.Endpoint)
	} else if dcfg.UseBoot2Docker {
		endpoint := os.Getenv("DOCKER_HOST")
		if endpoint != "" {
			path := os.Getenv("DOCKER_CERT_PATH")
			cert := fmt.Sprintf("%s/cert.pem", path)
			key := fmt.Sprintf("%s/key.pem", path)
			ca := fmt.Sprintf("%s/ca.pem", path)
			client, err = docker.NewTLSClient(endpoint, cert, key, ca)
		}
	} else {
		return Client{}, errors.New("empty docker configuration")
	}
	if err != nil || client == nil {
		return Client{dcfg, client}, err
	}

	return Client{dcfg, client}, client.Ping()
}

// IsValid returns true if the client is connected
func (c Client) IsValid() bool {
	return c.c != nil && c.c.Ping() == nil
}

// InspectImage returns the image stats
func (c Client) InspectImage(id ImageID) (Image, error) {
	img, err := c.c.InspectImage(string(id))
	ret := Image{ID: ImageID(img.ID)}
	if err != nil {
		return ret, err
	}
	if img.Config != nil {
		ret.Labels = img.Config.Labels
	}
	return ret, err
}

// ListImages lists the docker images
func (c Client) ListImages() ([]ImageID, error) {
	imgs, err := c.c.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		return nil, err
	}
	ret := make([]ImageID, 0, 0)
	for _, img := range imgs {
		ret = append(ret, ImageID(img.ID))
	}
	return ret, nil
}

// BuildImage builds a Docker image from a directory with a Dockerfile
func (c *Client) BuildImage(dirpath string) (Image, error) {
	var buf bytes.Buffer
	err := c.c.BuildImage(docker.BuildImageOptions{
		Dockerfile:   "Dockerfile",
		ContextDir:   dirpath,
		OutputStream: &buf,
	})
	var img Image
	if err != nil {
		return img, err
	}
	stepPrefix := "Step "
	successPrefix := "Successfully built "
	for err == nil {
		var line string
		line, err = buf.ReadString('\n')
		// fmt.Printf("build line: `%s`", line)
		if strings.HasPrefix(line, stepPrefix) {
			// step++
		} else if strings.HasPrefix(line, successPrefix) {
			img.ID = ImageID(strings.TrimSpace(line[len(successPrefix):]))
		}
	}
	err = c.c.Ping()
	if err != nil {
		log.Println("Warning: docker client lost after building image: ", err)
		var nc Client
		nc, err = NewClient(c.cfg)
		if err != nil {
			return img, err
		}
		c.c = nc.c
	}
	if err != nil && img.ID == "" {
		err = errors.New("unknown docker failure")
	}
	return c.InspectImage(img.ID)
}

// ExecuteImage takes a docker image, creates a container and executes it
func (c Client) ExecuteImage(id ImageID) (Container, error) {
	cfg := docker.Config{
		Image: string(id),
	}
	hc := docker.HostConfig{}
	cco := docker.CreateContainerOptions{
		Name:       "",
		Config:     &cfg,
		HostConfig: &hc,
	}
	cont, err := c.c.CreateContainer(cco)
	if err != nil {
		return Container{}, err
	}

	err = c.c.StartContainer(cont.ID, &hc)
	return Container{ID: ContainerID(cont.ID)}, err
}
