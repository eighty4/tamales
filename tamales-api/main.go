package main

import (
	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

var session *gocql.Session

func main() {
	cassandraHost := os.Getenv("CASSANDRA_HOST")
	if len(cassandraHost) == 0 {
		cassandraHost = "127.0.0.1"
	}
	cluster := gocql.NewCluster(cassandraHost)
	cluster.Keyspace = "tamales"
	session, _ = cluster.CreateSession()
	defer session.Close()

	router := httprouter.New()

	router.POST("/invite/vendor", inviteVendor)
	router.GET("/invite/vendor/:user_id", createVendor)

	router.POST("/login/initiate", requestLogin)
	router.GET("/login/validate/:login_token", login)

	router.POST("/vendor/:user_id/location", updateVendorLocation)
	router.GET("/vendors", getVendorLocations)
	router.GET("/vendor/:user_id/history", getVendorLocationHistory)

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
