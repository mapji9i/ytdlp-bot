package data

import (
	"database/sql"
	"fmt"
	"log"
	"ytdlp-bot/internal/environment"

	_ "github.com/mattn/go-sqlite3"
)

type RegisteredProcess struct {
	PID  int
	Date string
}

func ClearProcess() {
	db, err := sql.Open("sqlite3", environment.Environment.DBFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	db.Exec("DELETE FROM PIDS")
}

func PutProcessToDB(regProc RegisteredProcess) {

	db, err := sql.Open("sqlite3", environment.Environment.DBFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	_, err = db.Exec("insert into PIDS (PID, DATE_TIME_OF_CREATION) values ($1, $2)", regProc.PID, regProc.Date)
	if err != nil {
		log.Panic(err)
	}
}

func RemoveNextProcessNotEqualArgsFromDB(pid int) (RegisteredProcess, error) {
	db, err := sql.Open("sqlite3", environment.Environment.DBFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		log.Panic(err)
	}
	defer tx.Rollback()

	row := db.QueryRow(fmt.Sprintf("SELECT * FROM PIDS WHERE NOT PID=%d LIMIT 1", pid))
	regProc := RegisteredProcess{}
	err = row.Scan(&regProc.PID, &regProc.Date)
	if err != nil {
		return regProc, err
	}

	db.Exec("DELETE FROM PIDS WHERE NOT PID=$1", pid)

	tx.Commit()

	return regProc, nil
}
