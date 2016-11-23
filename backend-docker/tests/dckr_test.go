package tests

import (
	"log"
	"testing"
	"github.com/EUDAT-GEF/GEF/backend-docker/dckr"
	"github.com/EUDAT-GEF/GEF/backend-docker/config"
)

var configFilePath = "../config/config.json"
var clientConf []dckr.Config

func TestClient(t *testing.T) {
	settings, err := config.ReadConfigFile(configFilePath)
	if err != nil {
		t.Error("FATAL while reading config files: ", err)
	}
	clientConf = settings.Docker
	c := newClient(t)

	before := listImages(c, t)

	img := buildImage(c, t)

	after := listImages(c, t)

	errstr := "Cannot find new image in list"
	for _, x := range after {
		if x.ID == img.ID {
			errstr = ""
			break
		}
	}

	if errstr != "" {
		t.Error("before is: ", len(before), before)
		t.Error("image is: ", img)
		t.Error("after is: ", len(after), after)
		t.Error("")
		t.Error(errstr)
		t.Fail()
		return
	}

	containerID := executeImage(c, img.ID, t)
	log.Println("executed container: ", containerID)

	containerList := listContainers(c, t)
	if len(containerList) == 0 {
		t.Error("cannot find any containers")
		t.Fail()
		return
	}

	found := false
	for _, container := range containerList {
		if container.ID == containerID {
			found = true
		}
	}
	if !found {
		t.Error("cannot find the executed container in the list of all containers")
		t.Fail()
		return
	}

	inspectContainer(c, containerID, t)
}

func newClient(t *testing.T) dckr.Client {
	c, err := dckr.NewClientFirstOf(clientConf)
	if err != nil {
		t.Error(err)
		t.Error("--- client is not valid (this test requires a docker server)")
		t.Fail()
	}

	if !c.IsValid() {
		t.Error("client not valid (unable to ping)")
	}
	return c
}

func listImages(client dckr.Client, t *testing.T) []dckr.Image {
	imgs, err := client.ListImages()
	if err != nil {
		t.Error("List Image Error: ", err)
		t.Fail()
	}
	return imgs
}

func buildImage(client dckr.Client, t *testing.T) dckr.Image {
	img, err := client.BuildImage("./docker_test")
	if err != nil {
		t.Error("build image failed: ", err)
		t.Fail()
	}
	log.Println("built image:", img)
	return img
}

func executeImage(client dckr.Client, imgid dckr.ImageID, t *testing.T) dckr.ContainerID {
	containerID, err := client.ExecuteImage(imgid, nil)
	if err != nil {
		t.Error("starting image failed: ", err)
		t.Fail()
	}
	log.Println("starting image success: ", imgid)
	return containerID
}

func listContainers(client dckr.Client, t *testing.T) []dckr.Container {
	containers, err := client.ListContainers()
	if err != nil {
		t.Error("list containers failed: ", err)
		t.Fail()
	}
	log.Println("list containers success: ", containers)
	return containers
}

func inspectContainer(client dckr.Client, contID dckr.ContainerID, t *testing.T) dckr.Container {
	cont, err := client.InspectContainer(contID)
	if err != nil {
		t.Error("inspect container failed: ", err)
		t.Fail()
	}
	log.Println("inspect container success: ", cont)
	return cont
}
