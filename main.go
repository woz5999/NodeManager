package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/woz5999/NodeManager/pkg"
	config "github.com/woz5999/NodeManager/pkg/config"
	healthz "github.com/woz5999/NodeManager/pkg/healthz"
)

func main() {
	// Application configuration
	config, err := config.GetConfig()
	if err != nil {
		fmt.Printf("Failed to retrieve configuration: %s", err)
		os.Exit(1)
	}

	// Set the logging level.
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	// Instantiate the base struct.
	base, err := nodeman.NewBase(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Instantiate the application
	nodeman, err := nodeman.NewNodeMan(base)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Invoke the nodeman
	nodeman.Watch()

	// Health check.
	http.HandleFunc("/healthz", healthz.HandleFunc)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
