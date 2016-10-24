package dckr

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	"io/ioutil"
	"path/filepath"
)

const (
	minimalDockerVersion = 1006 // major * 1000 + minor
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

// VolumeID is a type for docker volume ids
type VolumeID string

//VolumeMountpoint is a type for docker volume mountpoint
type VolumeMountpoint string

// Image is a struct for Docker images
type Image struct {
	ID      ImageID
	RepoTag string
	Labels  map[string]string
}

// Container is a struct for Docker containers
type Container struct {
	ID    ContainerID
	Image Image
	State docker.State
}

type Volume struct{
	ID VolumeID
	Mountpoint VolumeMountpoint

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
			version, err := checkForMinimalDockerVersion(client.c)
			if err != nil {
				buf.WriteString(fmt.Sprintf(
					"%s:\n\t%s\nReason:%s\n",
					"Docker server version check has failed",
					dcfg, err))
			} else {
				log.Println("Successfully created Docker client using config:", dcfg)
				log.Println("Docker server version:", version)
				return client, nil
			}
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

func checkForMinimalDockerVersion(c *docker.Client) (string, error) {
	env, err := c.Version()
	if err != nil {
		return "", err
	}
	m := env.Map()
	version := m["Version"]
	arr := strings.Split(version, ".")
	if len(arr) < 2 {
		return "", fmt.Errorf("unparsable version string: %s", version)
	}
	major, err := strconv.Atoi(arr[0])
	if err != nil {
		return "", fmt.Errorf("unparsable major version: %s", version)
	}
	minor, err := strconv.Atoi(arr[1])
	if err != nil {
		return "", fmt.Errorf("unparsable minor version: %s", version)
	}
	if major*1000+minor < minimalDockerVersion {
		return "", fmt.Errorf("unusably old Docker version: %s", version)
	}
	return version, nil
}

// IsValid returns true if the client is connected
func (c Client) IsValid() bool {
	return c.c != nil && c.c.Ping() == nil
}

func makeImage(img *docker.Image) Image {
	repoTag := ""
	if len(img.RepoTags) > 0 {
		repoTag = img.RepoTags[0]
	}
	var labels map[string]string
	if img.Config != nil {
		labels = img.Config.Labels
	}
	return Image{
		ID:      ImageID(img.ID),
		RepoTag: repoTag,
		Labels:  labels,
	}
}

func makeImage2(img docker.APIImages) Image {
	repoTag := ""
	if len(img.RepoTags) > 0 {
		repoTag = img.RepoTags[0]
	}
	return Image{
		ID:      ImageID(img.ID),
		RepoTag: repoTag,
		Labels:  img.Labels,
	}
}

// InspectImage returns the image stats
func (c Client) InspectImage(id ImageID) (Image, error) {
	img, err := c.c.InspectImage(string(id))
	if err != nil {
		return Image{}, err
	}
	return makeImage(img), err
}

// ListImages lists the docker images
func (c Client) ListImages() ([]Image, error) {
	imgs, err := c.c.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		return nil, err
	}
	ret := make([]Image, 0, 0)
	for _, img := range imgs {
		ret = append(ret, makeImage2(img))
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
func (c Client) ExecuteImage(id ImageID, binds []string) (ContainerID, error) {
	img, err := c.c.InspectImage(string(id))
	log.Println("img", img)
	if err != nil {
		return ContainerID(""), err
	}
	hc := docker.HostConfig{
		Binds: binds,
	}
	log.Println("Config", img.Config);
	log.Println("binds", hc.Binds);
	cco := docker.CreateContainerOptions{
		Name:       "",
		Config:     &docker.Config{Image: img.ID},

		HostConfig: &hc,
	}
	fmt.Println("container options", cco)
	cont, err := c.c.CreateContainer(cco)
	if err != nil {
		log.Println("Create contianer failed")
		return ContainerID(""), err
	}

	err = c.c.StartContainer(cont.ID, &hc)
	return ContainerID(cont.ID), err
}




// ListContainers lists the docker images
func (c Client) ListContainers() ([]Container, error) {
	conts, err := c.c.ListContainers(
		docker.ListContainersOptions{All: true})
	if err != nil {
		return nil, err
	}
	ret := make([]Container, 0, 0)
	for _, cont := range conts {
		img, _ := c.InspectImage(ImageID(cont.Image))
		ret = append(ret, Container{
			ID:    ContainerID(cont.ID),
			Image: img,
			State: docker.State{
				Status: cont.Status,
			},
		})
	}
	return ret, nil
}

// InspectContainer returns the container details
func (c Client) InspectContainer(id ContainerID) (Container, error) {
	cont, err := c.c.InspectContainer(string(id))
	img, _ := c.InspectImage(ImageID(cont.Image))
	ret := Container{
		ID:    ContainerID(cont.ID),
		Image: img,
		State: cont.State,
	}
	if err != nil {
		return ret, err
	}
	return ret, err
}

// ListVolumes list all named volumes
func (c Client) ListVolumes() ([]Volume, error) {
	vols, err := c.c.ListVolumes(docker.ListVolumesOptions{})
	if err != nil {
		return nil, err
	}

	ret := make([]Volume, 0, 0)

	for _, vol := range vols{
		volume, _ := c.InspectVolume(VolumeID(vol.Name))
		ret = append(ret, Volume{
			ID:    VolumeID(volume.ID),
			Mountpoint: VolumeMountpoint(volume.Mountpoint),
		})
	}
	return ret, nil

}

// InspectVolume returns the volume details
func (c Client) InspectVolume(id VolumeID) (Volume, error) {
	volume, err := c.c.InspectVolume(string(id));
	ret := Volume{
		ID: VolumeID(volume.Name),
		Mountpoint: VolumeMountpoint(volume.Mountpoint),
	}
	return ret, err
}

// return file names in a directory
// TODO
func (c Client) ListVolumeContent(id VolumeID) []string {
	return nil
}

// BuildVolume creates a volume, copies data from dirpath
func (c Client) BuildVolume(dirpath string) (Volume, error) {
	volume, err := c.c.CreateVolume(docker.CreateVolumeOptions{});
	if err != nil {
		return Volume{}, err
	}
	ret := Volume{
		ID: VolumeID(volume.Name),
		Mountpoint: VolumeMountpoint(volume.Mountpoint),
	}
	//copy all content to volume
	err = copyDataToVolume(dirpath, string(ret.Mountpoint))
	return ret, err
}

func copyDataToVolume(src string, dst string) error{
	log.Printf("Copying data from %s to volume %s \n", src, dst)
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	copyTreeOptions := &CopyTreeOptions{Symlinks:true}


	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		srcFileInfo, err := os.Stat(srcPath)
		if err != nil {
			return err
		}

		if !srcFileInfo.IsDir() {
			//no dir
			err = CopyFile(srcPath, dstPath, false);
			if err != nil {
				log.Println("Copy files failed", err)
				return err
			}
		} else {
			//copy dir
			err = CopyTree(srcPath, dstPath, copyTreeOptions)
			if err != nil {
				log.Println("Copy directories failed", err)
				return err
			}
		}
	}
	return nil
}



//RemoveVolume removes a volume
func (c Client) RemoveVolume(id VolumeID) error {
	err := c.c.RemoveVolume(string(id))
	return err
}
