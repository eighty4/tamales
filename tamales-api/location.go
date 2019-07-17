package main

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

const saveVendorLocationCql = "insert into vendor_location (user_id, location_time, location) values (?, ?, ?);"
const saveVendorLocationHistoryCql = "insert into vendor_location_history (user_id, location_time, location) values (?, ?, ?);"

func updateVendorLocation(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var location VendorLocation
	if err := json.NewDecoder(request.Body).Decode(&location); err != nil {
		panic(err)
	}

	log.Printf("%s (%s): %s", location.VendorId, location.UpdateTime, *location.Location)

	batch := session.NewBatch(gocql.LoggedBatch)
	batch.Query(saveVendorLocationCql, location.VendorId, location.UpdateTime, location.Location)
	batch.Query(saveVendorLocationHistoryCql, location.VendorId, location.UpdateTime, location.Location)
	err := session.ExecuteBatch(batch)
	if err != nil {
		panic(err)
	}

	writer.WriteHeader(201)
}

const selectVendorsCql = "select user_id, location_time, location from vendor_location"

func getVendorLocations(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var locations []*VendorLocation
	var vendorId *gocql.UUID
	var updateTime *time.Time
	var location *string

	iter := session.Query(selectVendorsCql).Iter()
	for iter.Scan(&vendorId, &updateTime, &location) {
		vendorLocation := &VendorLocation{vendorId, updateTime, location}
		locations = append(locations, vendorLocation)
	}

	writer.Header().Set("Content-Type", "application/json;charset=utf-8")

	if len(locations) == 0 {
		if _, err := writer.Write([]byte("[]")); err != nil {
			panic(err)
		}
	} else {
		if err := json.NewEncoder(writer).Encode(locations); err != nil {
			panic(err)
		}
	}
}

const selectVendorLocationHistoryCql = "select location_time, location from vendor_location_history where user_id = ?;"

func getVendorLocationHistory(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	vendorId, err := gocql.ParseUUID(params.ByName("user_id"))
	if err != nil {
		panic(err)
	}

	log.Printf("GET /vendor/%s/location", vendorId)

	var locations []*VendorLocation
	var updateTime *time.Time
	var location *string

	iter := session.Query(selectVendorLocationHistoryCql, vendorId).Iter()
	for iter.Scan(&vendorId, &updateTime, &location) {
		vendorLocation := &VendorLocation{&vendorId, updateTime, location}
		locations = append(locations, vendorLocation)
	}

	writer.Header().Set("Content-Type", "application/json;charset=utf-8")

	if len(locations) == 0 {
		if _, err := writer.Write([]byte("[]")); err != nil {
			panic(err)
		}
	} else {
		if err := json.NewEncoder(writer).Encode(locations); err != nil {
			panic(err)
		}
	}
}
