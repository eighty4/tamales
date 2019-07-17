package main

import (
	"encoding/json"
	"github.com/eighty4/sse"
	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type InitiateLoginRequest struct {
	Email string `json:"email"`
}

type LoginResponse struct {
	LoginToken *gocql.UUID `json:"loginToken"`
	Expires    string      `json:"expires"`
}

const selectUserIdByEmailCql = "select user_id from user_lookup_by_email where email = ?;"
const saveLoginTokenCql = "insert into login_token (login_token, user_id) values (?, ?);"

var pendingLoginTokens = make(map[gocql.UUID]chan gocql.UUID)

func requestLogin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var initiateLoginRequest *InitiateLoginRequest
	if err := json.NewDecoder(request.Body).Decode(&initiateLoginRequest); err != nil {
		log.Printf("failed decoding request json: %s", err)
		http.Error(writer, "error initiating login", 500)
		return
	}

	var userId *gocql.UUID
	if err := session.Query(selectUserIdByEmailCql, initiateLoginRequest.Email).Scan(&userId); err != nil {
		log.Printf("failed looking up user by email: %s", err)
		http.Error(writer, "error initiating login", 500)
		return
	}

	loginToken := gocql.TimeUUID()

	if err := session.Query(saveLoginTokenCql, loginToken, userId).Exec(); err != nil {
		log.Printf("failed saving login token: %s", err)
		http.Error(writer, "error initiating login", 500)
		return
	}

	connection, err := sse.Upgrade(writer, request)
	if err != nil {
		log.Printf("failed upgrading connection to sse: %s", err)
		http.Error(writer, "error initiating login", 500)
		return
	}

	// todo send login email
	log.Printf("http://localhost:8080/login/validate/%s", loginToken)

	pendingLoginTokens[loginToken] = make(chan gocql.UUID)
	authToken := <-pendingLoginTokens[loginToken]
	if err := connection.SendJson(&LoginResponse{&authToken, "expires"}); err != nil {
		log.Printf("failed sending auth token as sse: %s", err)
		http.Error(writer, "error initiating login", 500)
		return
	} else {
		connection.Close()
		delete(pendingLoginTokens, loginToken)
	}
}

const selectLoginTokenCql = "select user_id from login_token where login_token = ?;"
const deleteLoginTokenCql = "delete from login_token where login_token = ? if exists;"
const insertAuthTokenCql = "insert into auth_token (auth_token, user_id) values (?, ?);"

func login(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	loginToken, err := gocql.ParseUUID(params.ByName("login_token"))
	if err != nil {
		log.Printf("failed parsing login token api param: %s", loginToken)
		http.ServeFile(writer, request, "static/vendor-login-error.html")
		return
	}

	var userIdString string
	selectVendorInviteIter := session.Query(selectLoginTokenCql, loginToken).Iter()
	if selectVendorInviteIter.NumRows() == 0 {
		log.Printf("failed to find login token: %s", loginToken)
		http.ServeFile(writer, request, "static/vendor-login-error.html")
		return
	} else {
		selectVendorInviteIter.Scan(userIdString)
	}

	userId, err := gocql.ParseUUID(params.ByName("login_token"))
	if err != nil {
		log.Printf("failed parsing uuid from string: %s", loginToken)
		http.ServeFile(writer, request, "static/vendor-login-error.html")
		return
	}

	if err := session.Query(deleteLoginTokenCql, loginToken).Exec(); err != nil {
		log.Printf("failed to delete login token: %s", err)
		http.ServeFile(writer, request, "static/vendor-login-error.html")
		return
	}

	authToken := gocql.TimeUUID()
	if err := session.Query(insertAuthTokenCql, authToken, userId).Exec(); err != nil {
		log.Printf("failed to save auth token: %s", err)
		http.ServeFile(writer, request, "static/vendor-login-error.html")
		return
	}

	pendingLoginTokens[loginToken] <- authToken
	http.ServeFile(writer, request, "static/vendor-login-success.html")
}
