package dckr

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	docker "github.com/fsouza/go-dockerclient"

	"github.com/EUDAT-GEF/GEF/backend-docker/def"
)

const (
	minimalDockerVersion = 1006 // major * 1000 + minor
)

// Client a Docker client with easy to use API
type Client struct {
	cfg def.DockerConfig
	c   *docker.Client
}

// ImageID is a type for docker image ids
type ImageID string

// ContainerID is a type for docker image ids
type ContainerID string

// VolumeID is a type for docker volume ids
type VolumeID string

// Image is a struct for Docker images
type Image struct {
	ID      ImageID
	RepoTag string
	Labels  map[string]string
	Created time.Time
	Size    int64
}

// Container is a struct for Docker containers
type Container struct {
	ID     ContainerID
	Image  Image
	State  docker.State
	Mounts []docker.Mount
}

// Volume is a struct for Docker volumes
type Volume struct {
	ID VolumeID
}

// VolBind is a binding between a Docker volume and a mounting path
type VolBind struct {
	VolumeID   VolumeID
	MountPoint string
	IsReadOnly bool
}

// NewClientFirstOf returns a new docker client or an error
func NewClientFirstOf(cfg []def.DockerConfig) (Client, error) {
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
func NewClient(dcfg def.DockerConfig) (Client, error) {
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

// InspectImage returns the image stats
func (c Client) InspectImage(id ImageID) (Image, error) {
	img, err := c.c.InspectImage(string(id))
	if err != nil {
		return Image{}, err
	}
	repoTag := ""
	if len(img.RepoTags) > 0 {
		repoTag = img.RepoTags[0]
	}
	var labels map[string]string
	if img.Config != nil {
		labels = img.Config.Labels
	}
	return Image{
		ID:      stringToImageID(img.ID),
		RepoTag: repoTag,
		Labels:  labels,
		Created: img.Created,
		Size:    img.Size,
	}, nil
}

func stringToImageID(id string) ImageID {
	id = strings.TrimSpace(id)
	shaPrefix := "sha256:"
	if strings.HasPrefix(id, shaPrefix) {
		id = id[len(shaPrefix):]
	}
	return ImageID(id)
}

// ListImages lists the docker images
func (c Client) ListImages() ([]Image, error) {
	imgs, err := c.c.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		return nil, err
	}
	makeImage := func(img docker.APIImages) Image {
		repoTag := ""
		if len(img.RepoTags) > 0 {
			repoTag = img.RepoTags[0]
		}
		return Image{
			ID:      stringToImageID(img.ID),
			RepoTag: repoTag,
			Labels:  img.Labels,
			Created: time.Unix(img.Created, 0),
			Size:    img.Size,
		}
	}

	ret := make([]Image, 0, 0)
	for _, img := range imgs {
		ret = append(ret, makeImage(img))
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
		if strings.HasPrefix(line, stepPrefix) {
			// step++
		} else if strings.HasPrefix(line, successPrefix) {
			img.ID = stringToImageID(line[len(successPrefix):])
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

// StartImage takes a docker image, creates a container and starts it
func (c Client) StartImage(id ImageID, cmdArgs []string, binds []VolBind) (ContainerID, error) {
	if id == "" {
		return ContainerID(""), def.Err(nil, "Empty image id")
	}
	fmt.Println("inspect image ", id)
	img, err := c.c.InspectImage(string(id))
	if err != nil {
		return ContainerID(""), def.Err(err, "InspectImage failed")
	}

	bs := make([]string, len(binds), len(binds))
	for i, b := range binds {
		bs[i] = fmt.Sprintf("%s:%s", b.VolumeID, b.MountPoint)
		if b.IsReadOnly {
			bs[i] = fmt.Sprintf("%s:ro", bs[i])
		}
	}

	config := *img.Config
	for _, arg := range cmdArgs {
		config.Cmd = append(config.Cmd, arg)
	}

	hc := docker.HostConfig{
		Binds: bs,
	}

	cco := docker.CreateContainerOptions{
		Config:     &config,
		HostConfig: &hc,
	}

	cont, err := c.c.CreateContainer(cco)
	if err != nil {
		return ContainerID(""), def.Err(err, "CreateContainer failed")
	}

	err = c.c.StartContainer(cont.ID, &hc)
	if err != nil {
		c.c.RemoveContainer(docker.RemoveContainerOptions{ID: cont.ID, Force: true})
		return ContainerID(""), def.Err(err, "StartContainer failed")
	}

	return ContainerID(cont.ID), nil
}

// ExecuteImage takes a docker image, creates a container and executes it, and waits for it to end
func (c Client) ExecuteImage(id ImageID, cmdArgs []string, binds []VolBind, removeOnExit bool) (int, error) {
	containerID, err := c.StartImage(id, cmdArgs, binds)
	if err != nil {
		//return 0, err
		return 0, def.Err(err, "StartImage failed")
	}

	//return c.WaitContainer(containerID, removeOnExit)
	return c.WaitContainer(containerID, false)
}

// DeleteImage removes an image by ID
func (c Client) DeleteImage(id string) error {
	err := c.c.RemoveImage(id)
	return err
}

// StartExistingContainer starts an existing container
func (c Client) StartExistingContainer(contID string, binds []string) (ContainerID, error) {
	hc := docker.HostConfig{
		Binds: binds,
	}

	err := c.c.StartContainer(contID, &hc)
	if err != nil {
		c.c.RemoveContainer(docker.RemoveContainerOptions{ID: contID, Force: true})
		return ContainerID(""), err
	}
	return ContainerID(contID), nil
}

// WaitContainer takes a docker container and waits for its finish.
// It returns the exit code of the container.
func (c Client) WaitContainer(id ContainerID, removeOnExit bool) (int, error) {
	containerID := string(id)
	exitCode, err := c.c.WaitContainer(containerID)
	if removeOnExit {
		c.c.RemoveContainer(docker.RemoveContainerOptions{ID: containerID, Force: true})
	}
	return exitCode, err
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
		img, _ := c.InspectImage(stringToImageID(cont.Image))
		mounts := make([]docker.Mount, 0, 0)
		for _, cont := range cont.Mounts {
			mounts = append(mounts, docker.Mount{
				Name:        cont.Name,
				Source:      cont.Source,
				Destination: cont.Destination,
				Driver:      cont.Driver,
				RW:          cont.RW,
				Mode:        cont.Mode,
			})
		}
		ret = append(ret, Container{
			ID:    ContainerID(cont.ID),
			Image: img,
			State: docker.State{
				Status: cont.Status,
			},
			Mounts: mounts,
		})
	}
	return ret, nil
}

// InspectContainer returns the container details
func (c Client) InspectContainer(id ContainerID) (Container, error) {
	cont, err := c.c.InspectContainer(string(id))
	img, _ := c.InspectImage(stringToImageID(cont.Image))
	ret := Container{
		ID:     ContainerID(cont.ID),
		Image:  img,
		State:  cont.State,
		Mounts: cont.Mounts,
	}
	if err != nil {
		return ret, err
	}
	return ret, err
}

// NewVolume builds an empty Docker volume
func (c *Client) NewVolume() (Volume, error) {
	cvo := docker.CreateVolumeOptions{}
	v, err := c.c.CreateVolume(cvo)
	if err != nil {
		return Volume{}, err
	}
	return Volume{
		ID: VolumeID(v.Name),
	}, nil
}

// ListVolumes list all named volumes
func (c Client) ListVolumes() ([]Volume, error) {
	vols, err := c.c.ListVolumes(docker.ListVolumesOptions{})
	if err != nil {
		return nil, err
	}

	ret := make([]Volume, 0, 0)
	for _, vol := range vols {
		ret = append(ret, Volume{
			ID: VolumeID(vol.Name),
		})
	}
	return ret, nil
}

//RemoveVolume removes a volume
func (c Client) RemoveVolume(id VolumeID) error {
	return c.c.RemoveVolume(string(id))
}

// GetTarStream returns a reader with a tar stream of a file path in a container
func (c Client) GetTarStream(containerID, filePath string) (io.Reader, error) {
	var b bytes.Buffer

	opts := docker.DownloadFromContainerOptions{
		Path:         filePath,
		OutputStream: &b,
	}

	err := c.c.DownloadFromContainer(containerID, opts)
	if err != nil {
		log.Println(filePath + " has not been retrieved")
	}

	return bytes.NewReader(b.Bytes()), err
}

// UploadSingleFile exported
func (c Client) UploadSingleFile(containerID, filePath string, dstPath string) error {
	var b bytes.Buffer

	fileHandler, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Cannot open " + filePath + ": " + err.Error())
		return err
	}

	tw := tar.NewWriter(&b)
	header, err := tar.FileInfoHeader(fileHandler, "")

	err = tw.WriteHeader(header)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	_, err = tw.Write(contents)
	if err != nil {
		log.Println("write to a file" + err.Error())
		return err
	}

	opts := docker.UploadToContainerOptions{
		Path:        dstPath,
		InputStream: &b,
	}

	err = c.c.UploadToContainer(containerID, opts)
	return err
}