package dckr

import (
	"log"
	"testing"
)

func TestClient(t *testing.T) {
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
	}

	containerID := executeImage(c, img.ID, t)
	log.Println("started container: ", containerID)

	containers := listContainers(c, t)
	if len(containers) == 0 {
		t.Error("cannot find any containers")
		t.Fail()
	}

	inspectContainer(c, containers[0].ID, t)
}

func newClient(t *testing.T) Client {
	config := []Config{
		Config{Endpoint: "unix:///var/run/docker.sock"},
		Config{UseBoot2Docker: true},
	}
	c, err := NewClientFirstOf(config)
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

func listImages(client Client, t *testing.T) []Image {
	imgs, err := client.ListImages()
	if err != nil {
		t.Error("List Image Error: ", err)
		t.Fail()
	}
	return imgs
}

func buildImage(client Client, t *testing.T) Image {
	img, err := client.BuildImage("./docker_test")
	if err != nil {
		t.Error("build image failed: ", err)
		t.Fail()
	}
	log.Println("built image:", img)
	return img
}

func executeImage(client Client, imgid ImageID, t *testing.T) ContainerID {
	containerID, err := client.ExecuteImage(imgid, nil)
	if err != nil {
		t.Error("starting image failed: ", err)
		t.Fail()
	}
	log.Println("starting image success: ", imgid)
	return containerID
}

func listContainers(client Client, t *testing.T) []Container {
	containers, err := client.ListContainers()
	if err != nil {
		t.Error("list containers failed: ", err)
		t.Fail()
	}
	log.Println("list containers success: ", containers)
	return containers
}

func inspectContainer(client Client, contID ContainerID, t *testing.T) Container {
	cont, err := client.InspectContainer(contID)
	if err != nil {
		t.Error("inspect container failed: ", err)
		t.Fail()
	}
	log.Println("inspect container success: ", cont)
	return cont
}
