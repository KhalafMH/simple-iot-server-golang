package models

import (
	"fmt"
	"github.com/gocql/gocql"
	"time"
)

var cluster *gocql.ClusterConfig
var session *gocql.Session

func GetReadingsFromDatabase(sensorId string, yearMonth string) (ReadingSlice, error) {
	if session == nil {
		initSession()
	}

	resultsIterator := session.Query(
		"SELECT id, type, value, alert, timestamp FROM readings WHERE id = ? AND year_month = ?",
		sensorId,
		yearMonth,
	).Iter()
	readings := make([]Reading, 0)

	var id string
	var typeValue string
	var value float64
	var alert bool
	var timestamp time.Time

	for resultsIterator.Scan(&id, &typeValue, &value, &alert, &timestamp) {
		reading := Reading{
			Id:        id,
			Type:      typeValue,
			Value:     value,
			Alert:     alert,
			Timestamp: timestamp,
		}
		readings = append(readings, reading)
	}
	err := resultsIterator.Close()
	if err != nil {
		return nil, err
	}

	return readings, nil
}

func GetReadingsFromDatabaseSince(sensorId string, since time.Time) (ReadingSlice, error) {
	if session == nil {
		initSession()
	}

	yearMonths := getYearMonthsSince(since)
	var resultsIterator *gocql.Iter
	if sensorId == "" {
		resultsIterator = session.Query(
			"SELECT id, type, value, alert, timestamp FROM readings WHERE timestamp >= ? ALLOW FILTERING",
			since,
		).Iter()
	} else {
		resultsIterator = session.Query(
			"SELECT id, type, value, alert, timestamp FROM readings "+
				"WHERE id = ? AND year_month IN ? AND timestamp >= ?",
			sensorId,
			yearMonths,
			since,
		).Iter()
	}

	readings := make([]Reading, 0)

	var id string
	var typeValue string
	var value float64
	var alert bool
	var timestamp time.Time

	for resultsIterator.Scan(&id, &typeValue, &value, &alert, &timestamp) {
		reading := Reading{
			Id:        id,
			Type:      typeValue,
			Value:     value,
			Alert:     alert,
			Timestamp: timestamp,
		}
		readings = append(readings, reading)
	}
	err := resultsIterator.Close()
	if err != nil {
		return nil, err
	}

	return readings, nil
}

func AddReadingToDatabase(reading Reading) error {
	if session == nil {
		initSession()
	}

	err := session.Query(
		"INSERT INTO readings (id, year_month, timestamp, type, value, alert) VALUES (?, ?, ?, ?, ?, ?)",
		reading.Id,
		reading.Timestamp.Format("2006-01"),
		reading.Timestamp.UTC(),
		reading.Type,
		reading.Value,
		reading.Alert,
	).Exec()

	return err
}

func getYearMonthsSince(since time.Time) []string {
	since = since.UTC()
	yearMonths := make([]string, 0)
	now := time.Now().UTC()

	for i := since; i.Year() < now.Year() || (i.Year() == now.Year() && i.Month() <= now.Month()); {
		yearMonths = append(yearMonths, i.Format("2006-01"))
		if i.Month() == time.December {
			i = time.Date(i.Year()+1, time.January, i.Day(), 0, 0, 0, 0, time.UTC)
		} else {
			i = time.Date(i.Year(), i.Month()+1, i.Day(), 0, 0, 0, 0, time.UTC)
		}
	}

	return yearMonths
}

func initCluster() {
	cluster = gocql.NewCluster("localhost")
	cluster.Keyspace = "sensor_app"
	cluster.Consistency = gocql.LocalOne
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: "sensor_app", Password: "Y6p4b152J2fZ"}
}

func initSession() {
	if cluster == nil {
		initCluster()
	}
	createSession, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("Cannot create cassandra session: %v", err))
	}
	session = createSession
}
