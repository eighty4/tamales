package main

import (
	"encoding/json"
	"github.com/eighty4/sse"
	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

type VendorInvite struct {
	Email string `json:"email"`
}

type VendorLocation struct {
	VendorId   *gocql.UUID `json:"vendorId"`
	UpdateTime *time.Time  `json:"updateTime"`
	Location   *string     `json:"location"`
}

type InitiateLoginRequest struct {
	Email string `json:"email"`
}

type LoginRequest struct {
	LoginToken *gocql.UUID `json:"loginToken"`
}

type LoginResponse struct {
	LoginToken *gocql.UUID `json:"loginToken"`
	Expires    string      `json:"expires"`
}

const saveVendorInviteCql = "insert into vendor_invite (vendor_id, email) values (?, ?);"
const saveVendorLookupByEmailCql = "insert into user_lookup_by_email (email, user_id, type) values (?, ?, 'vendor');"

func inviteVendor(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	vendorId, err := gocql.RandomUUID()
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	var invite VendorInvite
	if err := json.NewDecoder(request.Body).Decode(&invite); err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	batch := session.NewBatch(gocql.LoggedBatch)
	batch.Query(saveVendorInviteCql, vendorId, invite.Email)
	batch.Query(saveVendorLookupByEmailCql, invite.Email, vendorId)
	if err := session.ExecuteBatch(batch); err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	// todo send email with url
	log.Printf("http://localhost:8080/invite/vendor/%s", vendorId)
}

const selectVendorInviteCql = "select email from vendor_invite where vendor_id = ?;"
const deleteVendorInviteCql = "delete from vendor_invite where vendor_id = ?;"
const saveVendorCql = "insert into vendor (vendor_id, email) values (?, ?);"

func createVendor(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	vendorId := params.ByName("vendor_id")

	// select invite and validate
	selectVendorInviteIter := session.Query(selectVendorInviteCql, vendorId).Iter()
	if selectVendorInviteIter.NumRows() == 0 {
		writer.WriteHeader(401)
		http.ServeFile(writer, request, "static/vendor-invite-error.html")
		return
	}

	// get email from invite query
	var email *string
	selectVendorInviteIter.Scan(&email)

	// delete invite and save vendor
	batch := session.NewBatch(gocql.LoggedBatch)
	batch.Query(deleteVendorInviteCql, vendorId)
	batch.Query(saveVendorCql, vendorId, email)
	err := session.ExecuteBatch(batch)
	if err != nil {
		http.ServeFile(writer, request, "static/vendor-invite-error.html")
		return
	}

	http.ServeFile(writer, request, "static/vendor-getting-started.html")
}

const selectUserIdByEmailCql = "select user_id from user_lookup_by_email where email = ?;"
const saveLoginTokenCql = "insert into login_token (login_token, user_id) values (?, ?);"

var pendingLoginConnections = make(map[gocql.UUID]*sse.Connection)

func requestLogin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var initiateLoginRequest *InitiateLoginRequest
	if err := json.NewDecoder(request.Body).Decode(&initiateLoginRequest); err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	var userId *gocql.UUID
	if err := session.Query(selectUserIdByEmailCql, initiateLoginRequest.Email).Scan(userId); err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	loginToken := gocql.TimeUUID()

	if err := session.Query(saveLoginTokenCql, loginToken, userId).Exec(); err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	// todo send login email

	connection, err := sse.Upgrade(writer, request)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	pendingLoginConnections[loginToken] = connection

	// todo remove close connection from map on close
}

const selectLoginTokenCql = "select user_id from login_token where login_token = ?;"
const deleteLoginTokenCql = "delete from login_token where login_token = ? if exists;"

func login(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var loginRequest *LoginRequest
	if err := json.NewDecoder(request.Body).Decode(&loginRequest); err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	selectVendorInviteIter := session.Query(selectLoginTokenCql, loginRequest.LoginToken).Iter()
	if selectVendorInviteIter.NumRows() == 0 {
		http.Error(writer, "authentication failed", 401)
		return
	}

	if err := session.Query(deleteLoginTokenCql, loginRequest.LoginToken).Exec(); err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	loginSseConnection := pendingLoginConnections[*loginRequest.LoginToken]
	loginResponse := &LoginResponse{loginRequest.LoginToken, "asdf"}
	if err := loginSseConnection.SendJson(loginResponse); err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}
	loginSseConnection.Close()
}

const saveVendorLocationCql = "insert into vendor_location (vendor_id, location_time, location) values (?, ?, ?);"
const saveVendorLocationHistoryCql = "insert into vendor_location_history (vendor_id, location_time, location) values (?, ?, ?);"

func updateVendorLocation(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	log.Println("POST /vendor/location")

	var location VendorLocation
	if err := json.NewDecoder(request.Body).Decode(&location); err != nil {
		// todo handle error response
		panic(err)
	}

	log.Printf("%s (%s): %s", location.VendorId, location.UpdateTime, *location.Location)

	batch := session.NewBatch(gocql.LoggedBatch)

	batch.Query(saveVendorLocationCql, location.VendorId, location.UpdateTime, location.Location)
	batch.Query(saveVendorLocationHistoryCql, location.VendorId, location.UpdateTime, location.Location)

	err := session.ExecuteBatch(batch)
	if err != nil {
		// todo handle error response
		panic(err)
	}

	writer.WriteHeader(201)
}

const selectVendorsCql = "select vendor_id, location_time, location from vendor_location"

func getVendorLocations(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	log.Println("GET /vendor/location")

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
			// todo handle error response
			panic(err)
		}
	} else {
		if err := json.NewEncoder(writer).Encode(locations); err != nil {
			// todo handle error response
			panic(err)
		}
	}
}

const selectVendorLocationHistoryCql = "select location_time, location from vendor_location_history where vendor_id = ?;"

func getVendorLocationHistory(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	vendorId, err := gocql.ParseUUID(params.ByName("vendor_id"))
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
			// todo handle error response
			panic(err)
		}
	} else {
		if err := json.NewEncoder(writer).Encode(locations); err != nil {
			// todo handle error response
			panic(err)
		}
	}
}
