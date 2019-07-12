package main

import (
	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

var session *gocql.Session

func main() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "tamales"
	session, _ = cluster.CreateSession()
	defer session.Close()

	router := httprouter.New()

	router.POST("/invite/vendor", inviteVendor)
	router.GET("/invite/vendor/:vendor_id", createVendor)

	router.POST("/request-login", requestLogin)
	router.POST("/login", login)

	router.POST("/vendor/:vendor_id/location", updateVendorLocation)
	router.GET("/vendor", getVendorLocations)
	router.GET("/vendor/:vendor_id/history", getVendorLocationHistory)


	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
