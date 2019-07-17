package main

import (
	"encoding/json"
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

const saveVendorInviteCql = "insert into vendor_invite (user_id, email) values (?, ?);"
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

	// todo send invite email
	log.Printf("http://localhost:8080/invite/vendor/%s", vendorId)
}

const selectVendorInviteCql = "select email from vendor_invite where user_id = ?;"
const deleteVendorInviteCql = "delete from vendor_invite where user_id = ?;"
const saveVendorCql = "insert into vendor (user_id, email) values (?, ?);"

func createVendor(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	vendorId := params.ByName("user_id")

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
