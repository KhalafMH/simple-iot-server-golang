package main

import (
	"encoding/json"
	"fmt"
	"github.com/khalafmh/simple-iot-server-golang/models"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v <sensor_id>\n", os.Args[0])
		os.Exit(1)
	}

	sensorId := os.Args[1]

	for {
		value := rand.Float64()*50 - 25
		alert := false
		if value < -20 || value > 15 {
			alert = true
		}
		reading := models.Reading{Id: sensorId, Type: "temperature", Value: value, Alert: alert, Timestamp: time.Now().UTC()}
		marshal, err := json.Marshal(reading)
		if err != nil {
			continue
		}
		fmt.Printf("Sending reading: %v\n", string(marshal))
		_, err = http.Post(
			"http://localhost:8080/sensor-readings",
			"application/json",
			strings.NewReader(string(marshal)),
		)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		time.Sleep(1 * time.Second)
	}
}
