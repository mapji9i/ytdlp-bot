package data

import (
	"database/sql"
	"fmt"
	"log"

	"strings"
	"ytdlp-bot/internal/environment"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseEntity struct {
	VideoID       string
	VideoTitle    string
	ChannelTitle  string
	PlaylistTitle string
	InWork        int8
}

func CreateTable() {
	db, err := sql.Open("sqlite3", environment.Environment.DBFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Download_queue (video_id text PRIMARY KEY NOT NULL, VIDEO_TITLE TEXT NOT NULL, CHANNEL_TITLE TEXT NOT NULL, PLAYLIST_TITLE TEXT NOT NULL, IN_WORK INTEGER NOT NULL)")
	log.Print("Creating table")
	if err != nil {
		log.Print(err)
	}
}
func PutEntitiesToDB(entities []DatabaseEntity) {

	db, err := sql.Open("sqlite3", environment.Environment.DBFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	sb := strings.Builder{}
	sb.WriteString("insert into DOWNLOAD_QUEUE (VIDEO_ID, VIDEO_TITLE, CHANNEL_TITLE, PLAYLIST_TITLE, IN_WORK) values ")
	limit := len(entities) - 1
	for i, entity := range entities {
		sb.WriteString(fmt.Sprintf("('%s', '%s','%s', '%s',  %d)", entity.VideoID, entity.VideoTitle, entity.ChannelTitle, entity.PlaylistTitle, entity.InWork))
		if i < limit {
			sb.WriteString(", ")
		}
	}
	query := sb.String()
	db.Exec(query)
}

func PutEntityToDB(entity DatabaseEntity) {
	PutEntitiesToDB([]DatabaseEntity{entity})
}

func RemoveEntityFromDB(entity DatabaseEntity) {
	db, err := sql.Open("sqlite3", environment.Environment.DBFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	db.Exec("DELETE FROM DOWNLOAD_QUEUE WHERE VIDEO_ID=$1", entity.VideoID)
}
func ResetInWorkForAllEntities() {
	db, err := sql.Open("sqlite3", environment.Environment.DBFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	db.Exec("UPDATE DOWNLOAD_QUEUE SET IN_WORK=0 WHERE IN_WORK=1")
}
func ResetInWorkForEntity(entity DatabaseEntity) {
	db, err := sql.Open("sqlite3", environment.Environment.DBFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	db.Exec("UPDATE DOWNLOAD_QUEUE SET IN_WORK=0 WHERE VIDEO_ID=$1", entity.VideoID)
}
func GetNotUsedEntityInWork() (DatabaseEntity, error) {
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

	row := db.QueryRow("SELECT * FROM DOWNLOAD_QUEUE WHERE IN_WORK=0 LIMIT 1")
	entity := DatabaseEntity{}
	err = row.Scan(&entity.VideoID, &entity.VideoTitle, &entity.ChannelTitle, &entity.PlaylistTitle, &entity.InWork)
	if err != nil {
		return entity, err
	}
	entity.InWork = 1
	db.Exec("DELETE FROM DOWNLOAD_QUEUE WHERE VIDEO_ID=$1", entity.VideoID)
	db.Exec("INSERT INTO DOWNLOAD_QUEUE (VIDEO_ID, VIDEO_TITLE,  CHANNEL_TITLE, PLAYLIST_TITLE, IN_WORK) values ($1, $2, $3, $4, $5)", entity.VideoID, entity.VideoTitle, entity.ChannelTitle, entity.PlaylistTitle, entity.InWork)

	tx.Commit()

	return entity, nil
}

func ClearQueue() {
	db, err := sql.Open("sqlite3", environment.Environment.DBFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	db.Exec("DELETE FROM DOWNLOAD_QUEUE")
}
