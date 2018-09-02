package main

import (
	"fmt"
	"log"
	"net/http"
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
	log.Fatal(http.ListenAndServe(":3000", router))
}

func startRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	go recordFromCamera()
}

func stopRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	stopCamera()
}

func getRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
}

func recordFromCamera() {
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		fmt.Printf("Error on opening video capture.")
		return
	}
	defer webcam.Close()
	img := gocv.NewMat()
	defer img.Close()

	if ok := webcam.Read(&img); !ok {
		fmt.Printf("Cannot read camera device.")
		return
	}

	fileName := strconv.FormatInt(time.Now().Unix(), 10) + ".avi"
	writer, err := gocv.VideoWriterFile(fileName, "MJPG", 30, img.Cols(), img.Rows(), true)
	if err != nil {
		fmt.Printf("Error opening video writer device")
		return
	}
	defer writer.Close()
	startCamera()
	for recording > 0 {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed.")
			return
		}
		if img.Empty() {
			continue
		}
		writer.Write(img)
	}
}

func startCamera() {
	recording = 1
}

func stopCamera() {
	recording = 0
}
