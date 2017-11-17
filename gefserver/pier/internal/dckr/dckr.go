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

	"encoding/json"

	"context"

	"github.com/EUDAT-GEF/GEF/gefserver/def"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
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
	Cmd     []string
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

// VolumeInUse is an error that occurs when a volume is in use
var VolumeInUse = docker.ErrVolumeInUse

// NoSuchVolume is an error that occurs when a volume is not found
var NoSuchVolume = docker.ErrNoSuchVolume

// NodeAlreadyInSwarm is an error that occurs when we try to switch to the swarm mode but the node is already in a swarm
var NodeAlreadyInSwarm = docker.ErrNodeAlreadyInSwarm

// NodeNotInSwarm is an error that occurs when we try to switch to the swarm mode but the node is cannot do it
var NodeNotInSwarm = docker.ErrNodeNotInSwarm

// NewVolBind creates a new VolBind
func NewVolBind(id VolumeID, mount string, readonly bool) VolBind {
	return VolBind{
		VolumeID:   id,
		MountPoint: mount,
		IsReadOnly: readonly,
	}
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
	if !dcfg.TLSVerify {
		client, err = docker.NewClient(dcfg.Endpoint)
	} else {
		client, err = docker.NewTLSClient(dcfg.Endpoint, dcfg.CertPath, dcfg.KeyPath, dcfg.CAPath)
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
		Cmd:     img.Config.Cmd,
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

// LeaveIfInSwarmMode deactivates the Swarm Mode, if it was on
func (c Client) LeaveIfInSwarmMode() error {
	err := c.LeaveSwarmMode(true)
	if err == NodeNotInSwarm {
		return nil
	}
	return err
}

// InitiateSwarmMode switches a node to the Swarm Mode
func (c Client) InitiateSwarmMode(listenAddr string, advertiseAddr string) (string, error) {
	return c.c.InitSwarm(
		docker.InitSwarmOptions{
			InitRequest: swarm.InitRequest{
				ListenAddr:    listenAddr,
				AdvertiseAddr: advertiseAddr,
			},
		},
	)
}

// LeaveSwarmMode deactivates the Swarm Mode
func (c Client) LeaveSwarmMode(forced bool) error {
	return c.c.LeaveSwarm(docker.LeaveSwarmOptions{
		Force: forced,
	})
}

// GetSwarmContainerInfo finds a task associated with the given serviceID and retrieves information about the corresponding container
func (c Client) GetSwarmContainerInfo(serviceID string) (string, swarm.TaskState, error) {
	swarmTasks, err := c.c.ListTasks(docker.ListTasksOptions{})
	if err != nil {
		return "", swarm.TaskState(""), err
	}

	for _, task := range swarmTasks {
		if task.ServiceID == serviceID {
			err = nil
			if task.Status.Err != "" {
				err = def.Err(nil, task.Status.Err)
			}
			if task.Status.ContainerStatus.ContainerID != "" {
				return task.Status.ContainerStatus.ContainerID, task.Status.State, err
			}
			if task.Status.State == swarm.TaskStateComplete || task.Status.State == swarm.TaskStateFailed {
				return task.Status.ContainerStatus.ContainerID, task.Status.State, err
			}
		}
	}
	return "", swarm.TaskState(""), err
}

// isSwarmContainerStopped checks if a task container is not running
func (c Client) isSwarmContainerStopped(containerID string) (bool, swarm.Task, error) {
	task, err := c.findSwarmContainerTask(containerID)
	if err != nil {
		return false, swarm.Task{}, err
	}

	if task.Status.State == swarm.TaskStateComplete || task.Status.State == swarm.TaskStateFailed {
		var taskErr error
		if task.Status.Err != "" {
			taskErr = def.Err(nil, task.Status.Err)
		}
		return true, task, taskErr
	}

	return false, swarm.Task{}, nil
}

// findSwarmContainerStatus finds a swarm container status information by container ID
func (c Client) findSwarmContainerTask(containerID string) (swarm.Task, error) {
	swarmTasks, err := c.c.ListTasks(docker.ListTasksOptions{})
	if err != nil {
		return swarm.Task{}, err
	}

	for _, task := range swarmTasks {
		if task.Status.ContainerStatus.ContainerID == containerID {
			return task, nil
		}
	}
	return swarm.Task{}, def.Err(nil, "Could not find the container")
}

// WriteMonitor used to keep console output
type WriteMonitor struct{ io.Writer }

func (w *WriteMonitor) Write(bs []byte) (int, error) {
	n, err := w.Writer.Write(bs)
	log.Printf("Write() (%v, %v)", n, err)
	return n, err
}

// IsSwarmActive checks if the Swarm Mode is active
func (c Client) IsSwarmActive() (bool, error) {
	isActive := false
	dockerInfo, err := c.c.Info()
	if err != nil {
		return isActive, err
	}
	if dockerInfo.Swarm.LocalNodeState == "active" {
		isActive = true
	}
	return isActive, nil
}

// StartImage takes a docker image, creates a container and starts it
func (c Client) StartImage(id string, repoTag string, cmdArgs []string, binds []VolBind, limits def.LimitConfig, timeouts def.TimeoutConfig) (ContainerID, *bytes.Buffer, error) {
	var stdout bytes.Buffer

	if id == "" {
		return ContainerID(""), &stdout, def.Err(nil, "Empty image id")
	}

	if timeouts.Preparation == 0 {
		return ContainerID(""), &stdout, def.Err(nil, "Container preparation time out is not set")
	}

	if timeouts.JobExecution == 0 {
		return ContainerID(""), &stdout, def.Err(nil, "Job execution time out is not set")
	}

	img, err := c.c.InspectImage(string(id))
	if err != nil {
		return ContainerID(""), &stdout, def.Err(err, "InspectImage failed")
	}

	bs := make([]string, len(binds), len(binds))
	for i, b := range binds {
		bs[i] = fmt.Sprintf("%s:%s", b.VolumeID, b.MountPoint)
		if b.IsReadOnly {
			bs[i] = fmt.Sprintf("%s:ro", bs[i])
		}
	}

	config := *img.Config
	config.Cmd = cmdArgs

	config.AttachStdout = true
	config.AttachStderr = true
	hc := docker.HostConfig{
		Binds:      bs,
		CPUShares:  limits.CPUShares,
		CPUPeriod:  limits.CPUPeriod,
		CPUQuota:   limits.CPUQuota,
		Memory:     limits.Memory,
		MemorySwap: limits.MemorySwap,
	}

	createContainerContext, cancel := context.WithTimeout(context.Background(), time.Duration(timeouts.Preparation)*time.Second)
	defer cancel()

	cco := docker.CreateContainerOptions{
		Config:     &config,
		HostConfig: &hc,
		Context:    createContainerContext,
	}

	cont, err := c.c.CreateContainer(cco)
	if err != nil {
		return ContainerID(""), &stdout, def.Err(err, "CreateContainer failed")
	}

	attached := make(chan struct{})
	go func() {
		c.c.AttachToContainer(docker.AttachToContainerOptions{
			Container:    cont.ID,
			OutputStream: &stdout,
			ErrorStream:  &stdout,
			Logs:         true,
			Stdout:       true,
			Stderr:       true,
			Stream:       true,
			Success:      attached,
		})
	}()

	<-attached
	attached <- struct{}{}

	jobExecutionContext, cancel := context.WithTimeout(context.Background(), time.Duration(timeouts.JobExecution)*time.Second)
	defer cancel()
	err = c.c.StartContainerWithContext(cont.ID, &hc, jobExecutionContext)
	if err != nil {
		removeErr := c.TerminateContainerOrSwarmService(cont.ID, "")
		if removeErr != nil {
			log.Println(removeErr)
		}
		return ContainerID(""), &stdout, def.Err(err, "StartContainer failed")
	}

	return ContainerID(cont.ID), &stdout, nil
}

// StartSwarmService
func (c Client) StartSwarmService(id string, repoTag string, cmdArgs []string, binds []VolBind, limits def.LimitConfig, timeouts def.TimeoutConfig) (ContainerID, string, *bytes.Buffer, error) {
	var runningContainerID ContainerID

	swarmService, stdout, err := c.CreateSwarmService(repoTag, cmdArgs, binds, limits, timeouts)
	if err != nil {
		return runningContainerID, swarmService.ID, stdout, def.Err(err, "CreateSwarmService failed")
	}

	// Now we need to retrieve a container id
	for {
		contID, contState, err := c.GetSwarmContainerInfo(swarmService.ID)
		runningContainerID = ContainerID(contID)

		if (contID != "") || (contState == swarm.TaskStateComplete || contState == swarm.TaskStateFailed) {
			if err != nil {
				return runningContainerID, swarmService.ID, stdout, def.Err(err, "Failed to get information about the service container")
			}
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return runningContainerID, swarmService.ID, stdout, nil
}

// StartImageOrSwarmService
func (c Client) StartImageOrSwarmService(imgID string, imgRepoTag string, cmdArgs []string, binds []VolBind, limits def.LimitConfig, timeouts def.TimeoutConfig) (ContainerID, string, *bytes.Buffer, error) {
	var stdout *bytes.Buffer
	var runningContainer ContainerID
	var swarmService string

	swarmOn, err := c.IsSwarmActive()
	if err != nil {
		return runningContainer, swarmService, stdout, err
	}

	if swarmOn {
		runningContainer, swarmService, stdout, err = c.StartSwarmService(imgID, imgRepoTag, cmdArgs, binds, limits, timeouts)
		if err != nil {
			return runningContainer, swarmService, stdout, def.Err(err, "StartSwarmService failed")
		}
	} else {
		runningContainer, stdout, err = c.StartImage(imgID, imgRepoTag, cmdArgs, binds, limits, timeouts)
		if err != nil {
			return runningContainer, swarmService, stdout, def.Err(err, "StartImage failed")
		}
	}
	return runningContainer, swarmService, stdout, nil
}

// ExecuteImage takes a docker image, creates a container and executes it, and waits for it to end
func (c Client) ExecuteImage(imgID string, imgRepoTag string, cmdArgs []string, binds []VolBind, limits def.LimitConfig, timeouts def.TimeoutConfig, removeOnExit bool) (ContainerID, string, int, *bytes.Buffer, error) {
	runningContainer, swarmService, stdout, err := c.StartImageOrSwarmService(imgID, imgRepoTag, cmdArgs, binds, limits, timeouts)
	if err != nil {
		return runningContainer, swarmService, 0, stdout, err
	}

	exitCode, err := c.WaitContainerOrSwarmService(string(runningContainer))
	if err != nil {
		return runningContainer, swarmService, exitCode, stdout, def.Err(err, "WaitContainerOrSwarmService failed")
	}

	if removeOnExit {
		err = c.TerminateContainerOrSwarmService(string(runningContainer), swarmService)
	}
	return runningContainer, swarmService, exitCode, stdout, err
}

// DeleteImage removes an image by ID
func (c Client) DeleteImage(id string) error {
	err := c.c.RemoveImage(id)
	return err
}

// CreateSwarmService creates a Docker swarm service
func (c Client) CreateSwarmService(repoTag string, cmdArgs []string, binds []VolBind, limits def.LimitConfig, timeouts def.TimeoutConfig) (*swarm.Service, *bytes.Buffer, error) {
	var stdout bytes.Buffer
	var srv *swarm.Service

	if repoTag == "" {
		return srv, &stdout, def.Err(nil, "Empty image repoTag")
	}

	var serviceMounts []mount.Mount
	var curMount mount.Mount
	for _, v := range binds {
		curMount.ReadOnly = v.IsReadOnly
		curMount.Source = string(v.VolumeID)
		curMount.Target = v.MountPoint
		curMount.Type = mount.TypeVolume

		serviceMounts = append(serviceMounts, curMount)
	}

	/* Based on resources.CPUQuota = r.Limits.NanoCPUs * resources.CPUPeriod / 1e9
	taken from https://github.com/moby/moby/blob/v1.12.0-rc4/daemon/cluster/executor/container/container.go#L331 */
	calculatedNanoCPU := (limits.CPUQuota * 1e9) / limits.CPUPeriod

	swarmContext, cancel := context.WithTimeout(context.Background(), time.Duration(timeouts.JobExecution)*time.Second)
	defer cancel()

	serviceCreateOpts := docker.CreateServiceOptions{
		ServiceSpec: swarm.ServiceSpec{
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Image:   repoTag,
					Mounts:  serviceMounts,
					Command: cmdArgs,
				},
				RestartPolicy: &swarm.RestartPolicy{
					Condition: "none",
				},
				Resources: &swarm.ResourceRequirements{
					Limits: &swarm.Resources{
						NanoCPUs:    calculatedNanoCPU,
						MemoryBytes: limits.Memory,
					},
				},
			},
		},
		Context: swarmContext,
	}

	srv, err := c.c.CreateService(serviceCreateOpts)
	_, err = stdout.Write([]byte("/services/{id}/logs is an experimental feature introduced in Docker 1.13. Unfortunately, it is not yet supported by the Docker client we use"))
	if err != nil {
		return srv, &stdout, def.Err(err, "Failed to write a string into stdout stream")
	}

	return srv, &stdout, err
}

// waitSwarmService checks if a swarm service container is running till it is not (failed or done)
func (c Client) waitSwarmService(taskContainerId string) (int, error) {
	isStopped, task, err := c.isSwarmContainerStopped(taskContainerId)
	exitCode := task.Status.ContainerStatus.ExitCode
	for isStopped == false {
		isStopped, task, err = c.isSwarmContainerStopped(taskContainerId)
		if (task.Status.State == swarm.TaskStateComplete || task.Status.State == swarm.TaskStateFailed || task.Status.State == swarm.TaskStateShutdown) && (err != nil) {
			return 1, def.Err(err, "an error has occurred while executing a swarm service")
		}
		time.Sleep(200 * time.Millisecond)
	}
	return exitCode, err
}

// TerminateContainerOrSwarmService removes a container or a swarm service
func (c Client) TerminateContainerOrSwarmService(containerID string, swarmServiceID string) error {
	swarmOn, err := c.IsSwarmActive()
	if err != nil {
		return err
	}

	// Swarm mode (swarm services)
	if swarmOn {
		opts := docker.RemoveServiceOptions{ID: swarmServiceID}
		err = c.c.RemoveService(opts)
		if err != nil {
			if _, ok := err.(*docker.NoSuchService); ok {
				return nil
			}
			return def.Err(err, "an error has occurred while trying to remove a swarm service")
		}
		return nil
	} else {
		// Normal Mode (regular containers)
		removeOpts := docker.RemoveContainerOptions{
			ID:    containerID,
			Force: true,
		}

		err := c.c.RemoveContainer(removeOpts)
		if err != nil {
			if _, ok := err.(*docker.NoSuchContainer); ok {
				return nil
			}
			return def.Err(err, "an error has occurred while trying to remove a container")
		}
		return err
	}
}

// WaitContainerOrSwarmService takes a docker container/swarm service id, monitors it and waits for it to finish. An
// exit code of the container/swarm service is returned.
func (c Client) WaitContainerOrSwarmService(containerID string) (int, error) {
	swarmOn, err := c.IsSwarmActive()
	if err != nil {
		return 1, err
	}

	if swarmOn {
		// Swarm mode (swarm services)
		return c.waitSwarmService(containerID)
	} else {
		// Normal Mode (regular containers)
		return c.c.WaitContainer(containerID)
	}
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

// UploadFile2Container exported
func (c Client) UploadFile2Container(containerID, srcPath string, dstPath string) error {
	var b bytes.Buffer

	fileHandler, err := os.Stat(srcPath)
	if err != nil {
		log.Printf("Cannot open " + srcPath + ": " + err.Error())
		return err
	}

	tw := tar.NewWriter(&b)
	header, err := tar.FileInfoHeader(fileHandler, "")

	err = tw.WriteHeader(header)
	if err != nil {
		log.Println(def.Err(err, "UploadFile2Container/WriteHeader failed"))
		return err
	}

	contents, err := ioutil.ReadFile(srcPath)
	if err != nil {
		log.Println(def.Err(err, "UploadFile2Container/ReadFile failed"))
		return err
	}

	_, err = tw.Write(contents)
	if err != nil {
		log.Println("write to a file" + err.Error())
		return err
	}

	opts := docker.UploadToContainerOptions{
		Path:                 dstPath,
		InputStream:          &b,
		NoOverwriteDirNonDir: false,
	}

	err = c.c.UploadToContainer(containerID, opts)
	return err
}

func extractImageIDFromTar(imageFilePath string) (string, error) {
	type Manifest struct {
		Config   string
		RepoTags []string
		Layers   []string
	}

	var foundID string

	tarImage, err := os.Open(imageFilePath)
	if err != nil {
		return foundID, err
	}

	tarBallReader := tar.NewReader(tarImage)
	for {
		hdr, err := tarBallReader.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return foundID, def.Err(err, "reading tarball failed")
		}

		if strings.ToLower(hdr.Name) == "manifest.json" {
			var manifestContent []Manifest
			buf := new(bytes.Buffer)
			buf.ReadFrom(tarBallReader)

			err = json.Unmarshal(buf.Bytes(), &manifestContent)
			if err != nil {
				if err != nil {
					return foundID, def.Err(err, "cannot read manifest.json from the image")
				}
			}

			for _, value := range manifestContent {
				foundID := strings.Replace(value.Config, ".json", "", 1)
				if len(foundID) > 1 {
					return foundID, nil
				}
			}
		}
	}

	return foundID, def.Err(err, "could not retrieve image information")
}

// ImportImageFromTar installs a docker tar file as a docker image
func (c *Client) ImportImageFromTar(imageFilePath string) (ImageID, error) {
	var id string

	id, err := extractImageIDFromTar(imageFilePath)
	if err != nil {
		return ImageID(id), err
	}

	tar, err := os.Open(imageFilePath)
	if err != nil {
		return ImageID(id), err
	}
	defer tar.Close()

	opts := docker.LoadImageOptions{
		InputStream: tar,
	}

	err = c.c.LoadImage(opts)

	return ImageID(id), err
}

// TagImage tags a docker image
func (c *Client) TagImage(id string, repo string, tag string) error {
	opts := docker.TagImageOptions{
		Repo:  repo,
		Tag:   tag,
		Force: true,
	}

	return c.c.TagImage(id, opts)
}
