package main

import (
	"flag"
	"fmt"
	"net/http"
	"questionbasket/api"
	"questionbasket/config"
	"questionbasket/frame"
)

// for entire testing, must be removed after testing.
func main() {
	init := flag.Bool("init", false, "Run initial setup process")
	flag.Parse()

	config.LoadConfig()

	// Initialize Database
	dbConfig := frame.SQLDatabaseConnectionConfig{
		DBType:            frame.DRIVER_SQLITE,
		ConnectionAddress: config.Cfg.DatabaseURL,
	}
	frame.InitalizeDatabaseConnection(dbConfig)
	defer frame.CloseDatabaseConnection()

	if *init {
		DoInitSetup()
		return
	}

	APIServer := api.NewAPI()
	APIServer.RegisterPath()

	// Get MUX from router and register file server handler
	mux := APIServer.Router.GetMUX()
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs)

	//must be replaced to Certified lisnter
	fmt.Println("READY")
	http.ListenAndServe(":3000", mux)
}
