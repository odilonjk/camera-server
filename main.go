package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gocv.io/x/gocv"
)

var recording = 0

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/record/start", startRecord).Methods("POST")
	router.HandleFunc("/api/v1/record/stop", stopRecord).Methods("POST")
	router.HandleFunc("/api/v1/record", getRecord).Methods("GET")
	log.Println(http.ListenAndServe(":3000", router))
}

func startRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		returnError(w, 405)
	}
	go recordFromCamera()
}

func stopRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		returnError(w, 405)
	}
	stopCamera()
}

func getRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		returnError(w, 405)
	}
	lastVideo, err := getLastVideoName()
	if err != nil {
		log.Println(err)
		returnError(w, 500)
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+lastVideo)
	w.Header().Set("Content-Type", "video/x-msvideo")
	http.ServeFile(w, r, lastVideo)
}

func recordFromCamera() {
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Println("Error on opening video capture.")
		return
	}
	defer webcam.Close()
	img := gocv.NewMat()
	defer img.Close()

	if ok := webcam.Read(&img); !ok {
		log.Println("Cannot read camera device.")
		return
	}

	fileName := strconv.FormatInt(time.Now().Unix(), 10) + ".avi"
	writer, err := gocv.VideoWriterFile(fileName, "MJPG", 30, img.Cols(), img.Rows(), true)
	if err != nil {
		log.Println("Error opening video writer device")
		return
	}
	defer writer.Close()
	startCamera()
	for recording > 0 {
		if ok := webcam.Read(&img); !ok {
			log.Println("Device closed.")
			return
		}
		if img.Empty() {
			continue
		}
		writer.Write(img)
	}
}

func getLastVideoName() (string, error) {
	files := getFiles()
	var names []string
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".avi" {
			names = append(names, f.Name())
		}
	}
	if len(names) == 0 {
		return "", fmt.Errorf("there is no videos in the folder")
	}
	return names[len(names)-1], nil
}

func getFiles() []os.FileInfo {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Println(err)
	}
	return files
}

func startCamera() {
	recording = 1
}

func stopCamera() {
	recording = 0
}

func returnError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
	return
}
