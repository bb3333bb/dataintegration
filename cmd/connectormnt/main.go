package main

import (
	"log"

	"github.com/spf13/viper"
)

func main() {
	// Load configuration from file using Viper
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading configuration file: %s", err)
	}

	// Configure logging settings
	logFile := viper.GetString("log_file")
	logLevel := viper.GetString("log_level")
	if logFile == "" {
		log.Fatalf("Error: log_file not configured")
	}
	if logLevel == "" {
		logLevel = "info"
	}

	// Initialize logging to file and console
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logFile, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %s", err)
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	// Start application logic here
	log.Printf("Application starting...")
}