package main

import (
	"encoding/json"
	"github.com/khalafmh/simple-iot-server-golang/models"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/sensor-readings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			var resultReadings *models.ReadingSlice
			sinceQuery := r.URL.Query().Get("since")
			sensorIdQuery := r.URL.Query().Get("sensor_id")
			if sinceQuery != "" {
				since, err := time.Parse(time.RFC3339, sinceQuery)
				if err != nil {
					log.Printf("Error parsing `since` param as RFC3339 date. value: %v\n", sinceQuery)
					w.WriteHeader(400)
					return
				}
				res, err := models.GetReadingsFromDatabaseSince(sensorIdQuery, since)
				if err != nil {
					log.Printf("Error reading data from database: %v\n", err)
				}
				resultReadings = &res
			} else {
				now := time.Now().UTC()
				res, err := models.GetReadingsFromDatabaseSince(
					sensorIdQuery,
					time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC),
				)
				if err != nil {
					log.Printf("Error reading data from database: %v\n", err)
				}
				resultReadings = &res
			}
			b, err := json.Marshal(resultReadings)
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
				w.WriteHeader(500)
				return
			}
			var reading models.Reading
			err = json.Unmarshal(b, &reading)
			if err != nil {
				log.Println("Error decoding reading as JSON")
				w.WriteHeader(400)
				return
			}
			log.Printf("New sensor reading: %v\n", reading)
			err = models.AddReadingToDatabase(reading)
			if err != nil {
				log.Println("Error adding reading to database")
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(204)
		default:
			log.Printf("Unhandled request method %v\n", r.Method)
			w.WriteHeader(405)
		}
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
