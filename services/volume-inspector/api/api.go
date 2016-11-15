package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"time"
	"io/ioutil"
	"encoding/json"
	"os"
	"log"
)

// Volume folder content
type VolumeItem struct {
	Name       string `json:"name"`
	Size	   int64 `json:"size"`
	Modified   time.Time `json:"modified"`
	IsFolder   bool `json:"isFolder"`
	FolderTree []VolumeItem `json:"folderTree"`
}



type JReply struct {
	Message string `json:"message"`
}

type JPost struct {
	FolderPath string `json:"folderPath"`
}

const (
	jsonFileList = "/root/_filelist.json"
)

func readFolders(currentFolder string, volumeItems []VolumeItem) ([]VolumeItem, error) {
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
				}
			}
		}
	}

	return volumeItems, hasErrors
}

func Handlers() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", Index)
	router.HandleFunc("/ls", doLsRecursively).Methods("POST")
	log.Println("Starting server...")
	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	indexContent := JReply{Message:"Welcome to Volume Inspector"}
	json.NewEncoder(w).Encode(indexContent)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}

func doLsRecursively(w http.ResponseWriter, r *http.Request) {
	folderPath := ""
	// Form was POSTed
	if r.FormValue("folderPath") != "" {
		folderPath = r.FormValue("folderPath")
	} else { // JSON was POSTed
		var incomingData JPost
		if r.Body == nil {
			http.Error(w, "Please send a request body", http.StatusBadRequest)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&incomingData)
		if err != nil {
			http.Error(w, "Bad request: " + err.Error(), http.StatusBadRequest)
			return
		}
		folderPath = incomingData.FolderPath
	}

	if folderPath == "" {
		http.Error(w, "The path has not been specified", http.StatusBadRequest)
		return
	} else {
		doesExist, err := exists(folderPath)

		if doesExist {
			w.WriteHeader(http.StatusCreated)
			JFolderList, err := readFolders(folderPath, []VolumeItem{})
			if err == nil {
				json.NewEncoder(w).Encode(JFolderList)
				return
			} else {
				http.Error(w, "Bad request: " + err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "The folder you are trying to read does not exist", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, "Bad request: " + err.Error(), http.StatusBadRequest)
			return
		}
	}
}