package main

import (
	"encoding/json"
	"github.com/khalafmh/simple-iot-server-golang/models"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	readings := make([]models.Reading, 0, 10)
	http.HandleFunc("/sensor-readings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			b, err := json.Marshal(readings)
			if err != nil {
				log.Printf("Error: %v\n", err)
				return
			}
			_, err = w.Write(b)
			if err != nil {
				log.Printf("Error: %v\n", err)
				return
			}
		case "POST":
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading request body: %v", err)
				return
			}
			var reading models.Reading
			_ = json.Unmarshal(b, &reading)
			log.Printf("New sensor reading: %v\n", reading)
			readings = append(readings, reading)
		default:
			log.Printf("Unhandled request method %v\n", r.Method)
		}
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
