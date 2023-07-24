package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/bnallapeta/memcon/analysis"
	"github.com/bnallapeta/memcon/config"
	"github.com/bnallapeta/memcon/fetchProcesses"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var (
	appsEnvPath   string
	fetchInterval int
)

func init() {
	flag.StringVar(&appsEnvPath, "apps", "./config/apps.env", "Path to the .env file containing the list of apps to monitor")
	flag.IntVar(&fetchInterval, "interval", 360, "Interval in seconds for fetching memory consumption data")
	flag.Parse()
}

func setupDatabase(db *sql.DB) {
	statement, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS process_memory (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"timestamp" TEXT,
		"pid" INTEGER,
		"process_name" TEXT,
		"memory" INTEGER);`)
	statement.Exec()
}

func main() {

	config.InitConfig(appsEnvPath)
	db, err := sql.Open("sqlite3", "./process_memory.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	setupDatabase(db)

	go fetchProcesses.StoreProcesses(db, fetchInterval)

	router := gin.Default()

	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	router.GET("/visualize", func(c *gin.Context) {
		c.HTML(http.StatusOK, "visualize.html", nil)
	})

	router.GET("/visualize/data", func(c *gin.Context) {
		startTime := c.Query("start")
		endTime := c.Query("end")
		appName := c.Query("app")

		processes, err := analysis.FetchAppData(startTime, endTime, appName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Convert the timestamps to time.Time and sort the processes by timestamp in ascending order
		sort.SliceStable(processes, func(i, j int) bool {
			timestampI, err := time.Parse(time.RFC3339, processes[i].Timestamp)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return false
			}

			timestampJ, err := time.Parse(time.RFC3339, processes[j].Timestamp)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return false
			}

			return timestampI.Before(timestampJ)
		})

		c.JSON(http.StatusOK, processes)
	})

	router.Run(":8088")
}
