package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Apps to be monitored
var appsMonitored []string

// InitConfig initializes the configuration.
func InitConfig(envFile string) {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Error loading " + envFile + " file")
	}

	appsEnv := os.Getenv("APPS")
	if appsEnv == "" {
		log.Println("APPS environment variable is not set. Using default apps list.")
		appsEnv = "firefox,safari,brave,vscode,vscodium,terminal,gopls,windowserver,loginwindow"
	}

	appsMonitored = strings.Split(appsEnv, ",")
}

// GetApps returns the apps to be monitored.
func GetApps() []string {
	return appsMonitored
}
