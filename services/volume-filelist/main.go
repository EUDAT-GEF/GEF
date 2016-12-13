package main

import (
	"os"
	"encoding/json"
	"log"
	"time"
	"io/ioutil"
)

// Volume folder content
type VolumeItem struct {
	Name       string `json:"name"`
	Size	   int64 `json:"size"`
	Modified   time.Time `json:"modified"`
	IsFolder   bool `json:"isFolder"`
	FolderTree []VolumeItem `json:"folderTree"`
}

const (
	jsonFileList = "/root/_filelist.json"
	volumeFolder = "/root/volume"
)

func main() {
	jf, err := os.Create(jsonFileList)
	log.Println("Opening the JSON file")
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Reading the volume")
		JFolderList, err := readFolders(volumeFolder, []VolumeItem{})
		if err != nil {
			log.Println(err)
		} else {
			json.NewEncoder(jf).Encode(JFolderList)
			log.Println(JFolderList)
		}
	}
}

func readFolders(currentFolder string, volumeItems []VolumeItem) ([]VolumeItem, error) {
	log.Println("Reading folder: " + currentFolder)
	doesExist, hasErrors := exists(currentFolder)
	if hasErrors == nil {
		if doesExist {
			files, _ := ioutil.ReadDir(currentFolder)
			for _, f := range files {
				subFolderItems := []VolumeItem{}
				if f.IsDir() == true {
					subFolderItems, hasErrors = readFolders(currentFolder + "/" + f.Name(), []VolumeItem{})
				}
				if hasErrors == nil {
					volumeItems = append(volumeItems, VolumeItem{Name: f.Name(), Size: f.Size(), Modified: f.ModTime(), IsFolder:f.IsDir(), FolderTree: subFolderItems})
					log.Println(f.Name())
				}
			}
		}
	}

	return volumeItems, hasErrors
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}
