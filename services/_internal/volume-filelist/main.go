package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
	"path/filepath"
)

// Volume folder content
type VolumeItem struct {
	Name       string       `json:"name"`
	Size       int64        `json:"size"`
	Modified   time.Time    `json:"modified"`
	IsFolder   bool         `json:"isFolder"`
	Path   	   string       `json:"path"`
	FolderTree []VolumeItem `json:"folderTree"`
}

const (
	jsonFileList = "/root/_filelist.json"
	volumeFolder = "/root/volume"
)

func main() {
	subFolder := ""
	isRecursive := false
	if len(os.Args) > 1 {
		subFolder = os.Args[1]
	}
	if len(os.Args) > 2 {
		if os.Args[2] == "r" {
			isRecursive = true
		}
	}

	jf, err := os.Create(jsonFileList)
	log.Println("Opening the JSON file")
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Reading the volume")
		JFolderList, err := readFolders(subFolder, []VolumeItem{}, isRecursive)
		if err != nil {
			log.Println(err)
		} else {
			json.NewEncoder(jf).Encode(JFolderList)
			log.Println(JFolderList)
		}
	}
}

func readFolders(currentFolder string, volumeItems []VolumeItem, isRecursive bool) ([]VolumeItem, error) {
	log.Println("Reading folder: " + filepath.Join(volumeFolder, currentFolder))
	doesExist, hasErrors := exists(filepath.Join(volumeFolder, currentFolder))
	if hasErrors == nil {
		if doesExist {
			files, _ := ioutil.ReadDir(filepath.Join(volumeFolder, currentFolder))
			for _, f := range files {
				subFolderItems := []VolumeItem{}
				if f.IsDir() && isRecursive {
					subFolderItems, hasErrors = readFolders(filepath.Join(currentFolder, f.Name()), []VolumeItem{}, isRecursive)
				}
				if hasErrors == nil {
					volumeItems = append(volumeItems, VolumeItem{Name: f.Name(), Size: f.Size(), Modified: f.ModTime(), IsFolder: f.IsDir(),
						Path: currentFolder, FolderTree: subFolderItems})
					log.Println(f.Name())
				}
			}
		}
	}

	return volumeItems, hasErrors
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
