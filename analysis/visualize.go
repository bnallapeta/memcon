package analysis

import (
	"database/sql"
	"math"
	"strings"

	"github.com/bnallapeta/memcon/config"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ProcessAppsMemory struct {
	Timestamp     string
	Pid           uint64
	GroupName     string
	TotalMemory   uint64
	MaxMemory     uint64
	MinMemory     uint64
	AverageMemory float64
}

func FetchAppData(startTime string, endTime string, appName string) ([]ProcessAppsMemory, error) {
	db, err := sql.Open("sqlite3", "/Users/bnr/personal/memcon/process_memory.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT timestamp, pid, process_name, memory FROM process_memory 
	WHERE datetime(timestamp) BETWEEN datetime(?) AND datetime(?) AND process_name LIKE ? 
	ORDER BY datetime(timestamp);`, startTime, endTime, "%"+appName+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grouping := make(map[string]*ProcessAppsMemory)
	var totalCount int
	for rows.Next() {
		var timestamp, process_name string
		var pid, memory uint64
		err = rows.Scan(&timestamp, &pid, &process_name, &memory)
		if err != nil {
			return nil, err
		}

		groupName := getGroupName(process_name)

		// Convert memory to MB and round to 2 decimal places
		memory = uint64(math.Round(float64(memory)/1024/1024*100) / 100)

		if _, exists := grouping[groupName+timestamp]; !exists {
			grouping[groupName+timestamp] = &ProcessAppsMemory{Timestamp: timestamp, Pid: pid, GroupName: groupName, TotalMemory: memory, MaxMemory: memory, MinMemory: memory}
		} else {
			grouping[groupName+timestamp].TotalMemory += memory
			if memory > grouping[groupName+timestamp].MaxMemory {
				grouping[groupName+timestamp].MaxMemory = memory
			}
			if memory < grouping[groupName+timestamp].MinMemory {
				grouping[groupName+timestamp].MinMemory = memory
			}
		}
		totalCount += 1
	}

	var processList []ProcessAppsMemory
	for _, value := range grouping {
		value.AverageMemory = float64(value.TotalMemory) / float64(totalCount)
		processList = append(processList, *value)
	}

	return processList, nil
}

func getGroupName(process_name string) string {
	appsMonitored := config.GetApps()
	for _, app := range appsMonitored {
		if strings.Contains(strings.ToLower(process_name), strings.ToLower(app)) {
			return cases.Title(language.English).String(app)
		}
	}
	return "Other"
}
