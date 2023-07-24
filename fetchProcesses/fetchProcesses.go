package fetchProcesses

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bnallapeta/memcon/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shirou/gopsutil/process"
)

func StoreProcesses(db *sql.DB, fetchInterval int) {
	// Create table if it doesn't exist
	statement, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS process_memory (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"timestamp" TEXT,
			"pid" INTEGER,
			"process_name" TEXT,
			"memory" INTEGER);`)
	statement.Exec()

	appsMonitored := config.GetApps()
	for {
		var procCount = 0

		// Iterate over all running processes
		procs, err := process.Pids()
		if err != nil {
			log.Fatal(err)
		}

		for _, pid := range procs {
			proc, err := process.NewProcess(pid)
			if err != nil {
				if !strings.Contains(err.Error(), "process does not exist") {
					// Log the error but don't quit the program
					log.Printf("Could not create process object for pid %v: %v\n", pid, err)
				}
				continue
			}

			name, err := proc.Name()
			if err == nil {
				for _, procSubstring := range appsMonitored {
					if strings.Contains(strings.ToLower(name), strings.ToLower(procSubstring)) {
						storeProcInfo(db, proc)
						procCount = procCount + 1
						break
					}
				}
			}
		}

		fmt.Println("The number of processes inserted into the database are: ", procCount)
		// Sleep for fetchInterval seconds before the next iteration
		time.Sleep(time.Duration(fetchInterval) * time.Second)
	}
}

func storeProcInfo(db *sql.DB, proc *process.Process) {
	name, err := proc.Name()
	if err != nil {
		log.Printf("Could not fetch process name for pid %v: %v\n", proc.Pid, err)
		return
	}

	memInfo, err := proc.MemoryInfo()
	if err != nil {
		log.Printf("Could not fetch memory info for pid %v: %v\n", proc.Pid, err)
		return
	}

	timestamp := time.Now().Format(time.RFC3339)
	_, err = db.Exec("INSERT INTO process_memory (timestamp, pid, process_name, memory) VALUES (?, ?, ?, ?)",
		timestamp, proc.Pid, name, memInfo.RSS)
	if err != nil {
		log.Printf("Could not insert data into database: %v\n", err)
	}
}
